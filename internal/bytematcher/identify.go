// Copyright 2014 Richard Lehane. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bytematcher

import (
	"fmt"

	wac "github.com/richardlehane/match/fwac"
	"github.com/richardlehane/siegfried/pkg/config"
	"github.com/richardlehane/siegfried/pkg/core"
	"github.com/sniperkit/xfilter/internal/siegreader"
)

// identify function - brings a new matcher into existence
func (b *Matcher) identify(buf *siegreader.Buffer, quit chan struct{}, r chan core.Result, exclude ...int) {
	buf.Quit = quit
	waitSet := b.priorities.WaitSet(exclude...)
	var maxBOF, maxEOF int
	if len(exclude) > 0 {
		maxBOF, maxEOF = waitSet.MaxOffsets()
	} else {
		maxBOF, maxEOF = b.maxBOF, b.maxEOF
	}
	incoming := b.scorer(buf, waitSet, quit, r)
	rdr := siegreader.LimitReaderFrom(buf, maxBOF)
	// First test BOF frameset
	bfchan := b.bofFrames.index(buf, false, quit)
	for bf := range bfchan {
		if config.Debug() {
			fmt.Fprintln(config.Out(), strike{b.bofFrames.testTreeIndex[bf.idx], 0, bf.off, bf.length, false, true})
		}
		incoming <- strike{b.bofFrames.testTreeIndex[bf.idx], 0, bf.off, bf.length, false, true}
	}
	select {
	case <-quit: // the matcher has called quit
		for range bfchan {
		} // drain first
		close(incoming)
		return
	default:
	}

	// start bof matcher if not yet started
	b.bmu.Do(func() {
		b.bAho = wac.NewWac(b.lowmem, b.bofSeq.set)
	})
	var bchan chan wac.Result

	// Do an initial check of BOF sequences
	bchan = b.bAho.Index(rdr)
	for br := range bchan {
		if br.Index[0] == -1 {
			incoming <- progressStrike(br.Offset, false)
			if br.Offset > 131072 && (maxBOF < 0 || maxBOF > maxEOF*5) { // del buf.Stream 2^16	65536 2^17 131072
				break
			}
		} else {
			if config.Debug() {
				fmt.Fprintln(config.Out(), strike{b.bofSeq.testTreeIndex[br.Index[0]], br.Index[1], br.Offset, br.Length, false, false})
			}
			incoming <- strike{b.bofSeq.testTreeIndex[br.Index[0]], br.Index[1], br.Offset, br.Length, false, false}
		}
	}
	select {
	case <-quit: // the matcher has called quit
		for range bchan {
		} // drain first
		close(incoming)
		return
	default:
	}

	// Setup EOF tests
	efchan := b.eofFrames.index(buf, true, quit)
	b.emu.Do(func() {
		b.eAho = wac.NewWac(b.lowmem, b.eofSeq.set)
	})
	rrdr := siegreader.LimitReverseReaderFrom(buf, maxEOF)
	echan := b.eAho.Index(rrdr)

	// if we have a maximum value on EOF do a sequential search
	if maxEOF >= 0 {
		if maxEOF != 0 {
			_, _ = buf.CanSeek(0, true) // force a full read to enable EOF scan to proceed for streams
		}
		for ef := range efchan {
			if config.Debug() {
				fmt.Fprintln(config.Out(), strike{b.eofFrames.testTreeIndex[ef.idx], 0, ef.off, ef.length, true, true})
			}
			incoming <- strike{b.eofFrames.testTreeIndex[ef.idx], 0, ef.off, ef.length, true, true}
		}
		// Scan complete EOF
		for er := range echan {
			if er.Index[0] == -1 {
				incoming <- progressStrike(er.Offset, true)
			} else {
				if config.Debug() {
					fmt.Fprintln(config.Out(), strike{b.eofSeq.testTreeIndex[er.Index[0]], er.Index[1], er.Offset, er.Length, true, false})
				}
				incoming <- strike{b.eofSeq.testTreeIndex[er.Index[0]], er.Index[1], er.Offset, er.Length, true, false}
			}
		}
		// send a final progress strike with the maximum EOF
		incoming <- progressStrike(int64(maxEOF), true)
		// Finally, finish BOF scan
		for br := range bchan {
			if br.Index[0] == -1 {
				incoming <- progressStrike(br.Offset, false)
			} else {
				if config.Debug() {
					fmt.Fprintln(config.Out(), strike{b.bofSeq.testTreeIndex[br.Index[0]], br.Index[1], br.Offset, br.Length, false, false})
				}
				incoming <- strike{b.bofSeq.testTreeIndex[br.Index[0]], br.Index[1], br.Offset, br.Length, false, false}
			}
		}
		close(incoming)
		return
	}
	// If no maximum on EOF do a parallel search
	for {
		select {
		case br, ok := <-bchan:
			if !ok {
				if maxBOF < 0 && maxEOF != 0 {
					_, _ = buf.CanSeek(0, true) // if we've a limit BOF reader, force a full read to enable EOF scan to proceed for streams
				}
				bchan = nil
			} else {
				if br.Index[0] == -1 {
					incoming <- progressStrike(br.Offset, false)
				} else {
					if config.Debug() {
						fmt.Fprintln(config.Out(), strike{b.bofSeq.testTreeIndex[br.Index[0]], br.Index[1], br.Offset, br.Length, false, false})
					}
					incoming <- strike{b.bofSeq.testTreeIndex[br.Index[0]], br.Index[1], br.Offset, br.Length, false, false}
				}
			}
		case ef, ok := <-efchan:
			if !ok {
				efchan = nil
			} else {
				if config.Debug() {
					fmt.Fprintln(config.Out(), strike{b.eofFrames.testTreeIndex[ef.idx], 0, ef.off, ef.length, true, true})
				}
				incoming <- strike{b.eofFrames.testTreeIndex[ef.idx], 0, ef.off, ef.length, true, true}
			}
		case er, ok := <-echan:
			if !ok {
				echan = nil
			} else {
				if er.Index[0] == -1 {
					incoming <- progressStrike(er.Offset, true)
				} else {
					if config.Debug() {
						fmt.Fprintln(config.Out(), strike{b.eofSeq.testTreeIndex[er.Index[0]], er.Index[1], er.Offset, er.Length, true, false})
					}
					incoming <- strike{b.eofSeq.testTreeIndex[er.Index[0]], er.Index[1], er.Offset, er.Length, true, false}
				}
			}
		}
		if bchan == nil && efchan == nil && echan == nil {
			close(incoming)
			return
		}
	}
}
