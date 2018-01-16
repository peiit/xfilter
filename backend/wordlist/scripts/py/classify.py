# coding: utf-8
"""
This code is for chinese classifier use `Grocery`
And `Grocery` is based on `LibLinear`

Author: nicksite68@gmail.com
Time:   2016-03-31
"""

from tgrocery import Grocery
import csv, json

class MyGrocery(object):
  def __init__(self, name):
    super(MyGrocery, self).__init__()
    self.grocery = Grocery(name)
    self.loaded = False
    self.correct = 1.0

  def train(self, src):
    lines = []
    for line in csv.reader(open(src)):
      label, s = line[0],line[1]
      text = s.decode('utf8')
      lines.append((label, text))
    self.grocery.train(lines)

  def save_model(self):
    self.grocery.save()

  def train_and_save(self, src):
    self.train(src)
    self.save_model()

  def load_model(self):
    if not self.loaded:
      self.grocery.load()
      self.loaded = True

  def predict(self, text):
    self.load_model()
    return self.grocery.predict(text)

  def test(self, src):
    self.load_model()
    total, wrong_num = 0.0, 0.0
    for line in csv.reader(open(src)):
      total += 1
      if line[0] != self.predict(line[1]):
        wrong_num += 1

    print "load test file from " + src
    correct = (total - wrong_num ) / total
    self.correct = correct
    print "total: %d , wrong_num: %d, success percentage: %f" %(total, wrong_num, correct)
    result = dict(type="test", total=total, wrong_num=wrong_num, correct=correct)
    return json.dumps(result)

def main():
  grocery = MyGrocery("SVM")
  DATA_DIR = '../../data/pre/'
  grocery.train_and_save(DATA_DIR + 'train.csv')
  grocery.test(DATA_DIR + "test.csv")
  print grocery.predict("你好")

if __name__ == '__main__':
  main()