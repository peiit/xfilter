$(document).ready(function() {
  var postRequest = function (url, data, success_cb) {
    $.ajax({
      type:'post',
      url: url,
      data: JSON.stringify(data),
      dataType:'json',
      beforeSend:function(xhr){
        xhr.setRequestHeader('Content-Type', 'application/json' )
      },
      success:function(res){
        success_cb(res);
      },
      error:function(){
        console.log('failed!');
      }
    });
  }

  var uploadRequest = function (data, success_cb){
    $.ajax({
      type: 'post',
      url: 'upload',
      data: data,
      dataType:'json',
      success: function(res){
        if (res.status == "200") {
          success_cb(res)
        }else{
          alert("文件名或格式出错！");
        }
      },
      error: function() {
        console.log('failed!');
      },
      cache: false,
      contentType: false,
      processData: false
    })
  }

  var classifyTextRequest = function (data, success_cb){
    postRequest('classify', data, success_cb);
  }

  var ActionRequest = function (data, success_cb){
    postRequest('action', data, success_cb);
  }

  var handleClick = function(){
    var text = $("#input").val();
    var predictArea = $("#predict");
    if (text != "") {
      var data = {
        text: text
      };
      classifyTextRequest(data, function(data){
        var polyfill = [];
        polyfill.push("<div class='sample'>")
        polyfill.push("<p class='sample-filter'>过滤结果：", data.filter_result,"</p>");
        polyfill.push("<p class='sample-predict'>预测结果：", data.predict_result,"</p>");
        polyfill.push("</div>")
        predictArea.append(polyfill.join(""))
      });

    }else{
      alert("输入不能为空！")
    }
  };

  var handleUpload = function (e){
    e.preventDefault();
    var formData = new FormData($("form")[0]);
    var uploadButton = e.target;
    var message = $("#message");

    uploadRequest(formData, function(data){
      var filename = data.filename;
      var type = filename.split('.')[0];
      var typeMessage;
      switch(type){
        case "train":
          typeMessage = "训练";
          break;
        case "test":
          typeMessage = "测试";
          break;
        default:
          typeMessage = type
          break;
      }
      if (type == "sensitive") {
        alert("敏感词典上传成功！");
        window.setTimeout(function(){
          window.location.reload();
        },200)
        return;
      }
      var polyfill = [];
      var a = "<a id='action' data=" + type + ":" + filename + " href='javascript:void(0)'>"+ typeMessage +"</a>"
      polyfill.push("<p>", typeMessage, "文件上传成功！", "<span class='uploaded'>", "点此", a,"</span></p>");
      message.append(polyfill.join(""));
      uploadButton.value = "已上传";

      $("#action").bind("click", handleAction);
    })
  };

  var handleAction = function(e) {
    e.preventDefault();
    var elem = e.target;
    var actionType = elem.getAttribute("data").split(":");
    var data = {
      type: actionType[0],
      filename: actionType[1]
    }

    ActionRequest(data, function(data){
      var type = data.type;
      var message = $("#message");

      if (type=="train"){
        alert("模型训练成功！")
      }else{
        var polyfill = [];
        var correct = String(data.correct).slice(0, 8) * 100 + "%";
        polyfill.push("<p>测试结果: </p>", " 总数: ", data.total, " 错误数目: ", data.wrong_num, " 正确率: ", correct);
        message.append(polyfill.join(''));
      }
    })
  };


  $("#button").bind("click", handleClick);
  $("#upload").bind("click", handleUpload);
});