var dom;
var _this;
class CWarnInfo{
    constructor(style){
        this.type = style;
        _this = this;
     }
    initBack(){
        var link = document.createElement("link");
        link.rel = "stylesheet";
        link.type = "text/css";
        link.href = "../css/warning.css";
        var head = document.getElementsByTagName("head")[0];
        head.appendChild(link);
        if(!dom)
            dom = window.top.document.getElementsByTagName("body").item(0);
     }
     cancel_default(){
        
        var mask = dom.getElementsByClassName("_screen_mask").item(0);
        var content= dom.getElementsByClassName("_warn_content").item(0);
        if(mask && content){
            mask.remove();
            content.remove();
        }
        if(_this.cancel_call){
            _this.cancel_call();
        }
    }
  ok_default(){
    if(_this.ok_call){
        var value = "";
        if(_this.inputBtn){
            value= _this.inputBtn.value;
        }
        _this.ok_call(value);
    }
     _this.cancel_default();
 }

     show_warn(titleStr,contStr,ok_call,cancel_call){
         this.initBack();
         var mask = document.createElement("div");
         var content= document.createElement("div");
         var title = document.createElement("p");
         var info = document.createElement("p");
         if(this.type == "input"){
             info = document.createElement("input");
             info.placeholder=contStr;
             this.inputBtn = info;
         }
         var ok = document.createElement("button");
         var cancel = document.createElement("button");
         mask.className="_screen_mask";
         content.className="_warn_content";
         title.className="warn_head";
         info.className="warn_head warn_info";
         ok.className="warn_btn warn_btn_ok";
         cancel.className="warn_btn";
         if(titleStr){
             title.innerText = titleStr;
         }
         if(contStr && this.type!="input"){
             info.innerText = contStr;
         }
         ok.innerText="确定";
         cancel.innerText="取消";
         if(ok_call){
            this.ok_call = ok_call;
            console.log(this);

         }else{
             ok.style.width="60%";
             ok.innerText="我知道了";
         }

         if(cancel_call){
             this.cancel_call=cancel_call;
         }
          ok.addEventListener("click",this.ok_default);
          cancel.addEventListener("click",this.cancel_default);

         if(dom && mask){
             content.appendChild(title);
             content.appendChild(info);
             if(ok_call)
                 content.appendChild(cancel);
             content.appendChild(ok);
             dom.appendChild(mask);
             dom.appendChild(content);
         }
     }
}

function show_warn(titleStr,contStr,ok_call,cancel_call){
    var warn = new CWarnInfo("generic");
    _this=warn;
    warn.show_warn(titleStr,contStr,ok_call,cancel_call);
}
