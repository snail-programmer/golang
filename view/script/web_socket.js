class SockServ{
    constructor(callback){
        this.wsurl = "ws://127.0.0.1:9000/socknet.sock";
        this.sock = new WebSocket(this.wsurl);
        this.sock.onopen=function(){
            callback();
        }
    }
    sendData(data,callback){
       // this.sock = new WebSocket(this.wsurl);
        var _this = this;
        this.callback = callback;
        if(this.sock.readyState == 1){
            this.sock.send(data);
            console.log("send");
            _this.receive();
        }else{
            this.sock = new WebSocket(this.wsurl);
            this.sock.onopen=function(e){
                this.send(data);
                console.log("send2");
                _this.receive();
             }
        }
    }
    receive(){
        var _this = this;
        this.sock.onmessage = function(e){
            if(_this.callback){
                _this.callback(e.data);
            }
        }
        this.sock.onclose = function(e){
            console.log("close socket");
        }
    }
    
}
function util_connect(callback){
          //消息通知
          sock = new SockServ(function(e){
            sock.sendData(sessionStorage.getItem("MyUserId"),function(e){
            console.log("e.data:",e);
            sessionStorage.setItem("msg_cnt",e);
            if(callback){
                callback();
            }
        });
        });
}
 