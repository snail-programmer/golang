class CxhrServer{
     constructor(callback){
        this.xhrhttp=new XMLHttpRequest();
        this.callback = callback;
    }
    xhrsend(method,route,data){
        var xh = this.xhrhttp;
        var callback = this.callback;
        xh.onreadystatechange=function()
        {
            if(xh.readyState==4 && xh.status==200)
            {
                var response=xh.responseText;
        
                try {
                    var json1= JSON.parse(response);
                    if((typeof json1).toString() == "object"){
                       if(isDispatcherCmd(json1)){
                           return;
                       }
                    }
                    var xData=json1["data"];
                    console.log("xData:",  xData)
                     callback(xData);
                } catch (error) {
                    console.log("error:",error);
                    callback(response);
                }
            
            }
        }
        this.xhrhttp.open(method,route,true);
        this.xhrhttp.withCredentials=true; //允许发送cookies
        if(method == "POST"){
            this.xhrhttp.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
        }
        this.xhrhttp.send(data);
    }

};
//返回的是一个页面的话，转发
function isDispatcherCmd(obj){
     var jump = obj.dispatch;
     if(jump){
         document.location=jump;
         return true;
     }
     return false;
}
function request_route(route,callback){
    var xhr = new CxhrServer(callback);
    xhr.xhrsend('GET',route,"");
}
function send_data(route,callback,data){
    var xhr = new CxhrServer(callback);
    xhr.xhrsend('POST',route,data);
}