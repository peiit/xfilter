# coding: utf-8
from flask import Flask, request, jsonify, render_template
from werkzeug import secure_filename
from classify import MyGrocery
import logging, sys, os, requests, json

app = Flask(__name__)
# upload config
DATA_FOLDER = '../../data'
UPLOAD_FOLDER = DATA_FOLDER +'/uploads'
ALLOWED_EXTENSIONS = set(['csv','txt'])

# api host
API_HOST = 'http://127.0.0.1:8001/word/is_valid'

# log
app.logger.addHandler(logging.StreamHandler(sys.stdout))
app.logger.setLevel(logging.DEBUG)

# trained model init
grocery = MyGrocery("SVM")
grocery.predict("你好")

# utils
def request_params(request):
  if len(request.form) != 0:
    return request.form
  if len(request.json) != 0:
    return request.json

def allowed_file(filename, filetype=None):
  return '.' in filename and \
    filename.rsplit('.', 1)[1] in (filetype or ALLOWED_EXTENSIONS)

# routes
@app.route('/')
def index():
  return render_template('index.html')

@app.route('/upload', methods=['post'])
def upload():
  if request.method == 'POST':
    file = request.files['file']
    if file and allowed_file(file.filename):
      filename = secure_filename(file.filename)
      params = request_params(request)
      if params['type'] == 'sensitive' and allowed_file(filename, set(['txt'])) :
        filename = 'sensitive.txt'
        src = os.path.join(DATA_FOLDER, filename)
        file.save(src)
        res = { "status": 200, "filename": filename }
      elif params['type'] in ['test', 'train'] and allowed_file(filename, set(['csv'])):
        filename = params['type'] + '.csv'
        src = os.path.join(UPLOAD_FOLDER, filename)
        file.save(src)
        res = { "status": 200, "filename": filename }
      else:
        res = { "status": 500 }
    else:
      res = { "status": 500 }
  return jsonify(res)

@app.route('/action', methods=['POST'])
def action():
  params = request_params(request)
  action_type = params['type']
  filename = params['filename']
  src = UPLOAD_FOLDER + '/' + filename
  if action_type == "train":
    grocery.train_and_save(src)
    return jsonify({ "type" : 'train'})
  elif action_type == "test":
    result = grocery.test(src)
    return result
  else:
    return jsonify({ "status": 200 })

@app.route('/classify', methods=['POST'])
def classify():
  text = request_params(request)['text'].strip(' ')
  label = grocery.predict(text)
  res = requests.post(API_HOST, data=dict(v=text))
  print res.content
  app.logger.info("[INFO] label: "+ label + " text: " + text)
  predict_result = json.dumps(dict(label=label, text=text), ensure_ascii=False, encoding='utf8')
  return jsonify({ "predict_result": predict_result, "filter_result": res.content})

@app.route('/predict', methods=['POST'])
def predict():
  text = request_params(request)['text'].strip(' ')
  label = grocery.predict(text)
  app.logger.info("[INFO] label: "+ label + " text: " + text)
  return jsonify({ "label": label, "text": text})

if __name__ == '__main__':
  app.run(debug=True, port=8006, host='0.0.0.0')