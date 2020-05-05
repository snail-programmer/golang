request_route('/getMyUser',request_route_callback);
var globaldoc=window.top.document;
function request_route_callback(obj) {
    //隐藏注册登录按钮，显示用户图标
    var head=document.getElementById("loginhead");
    var nickname = document.getElementById("nickname");
    var btn=document.getElementById("btnStyle");
    if(obj.Id){
        localStorage.setItem("MyUserId",obj.Id);
        if(head && btn){
            head.style.visibility="visible";
            btn.style.visibility="hidden";
            head.style.background="url("+obj.HeadImg+") no-repeat";
            head.style.backgroundSize="cover";
            nickname.innerText = obj.NickName;
        }  
    }else{
        localStorage.removeItem("MyUserId");
        if(head && btn){
            head.style.visibility="hidden";
            btn.style.visibility="visible";
        }
    }
}
var jumpUrl = "";
var loadtitle = "";
function jumpcall(jump) {
    if(jump  == "未登录呢" && jumpUrl){
        show_maskTime(loadtitle,1000,function(){
            globaldoc.location=jumpUrl; 
        }); 
    }
 }
function check_identify_jump(title,url) {
    loadtitle = title;
    jumpUrl = url;
    request_route('/getMyUser',jumpcall);
}
function loginClick() {
    globaldoc.location = "login.html";
}
function reCheck(){
    request_route('/getMyUser',request_route_callback);
}
// function rootCheck(){
//     //获取顶层窗体
//     var dom = window.top.document.getElementsByTagName("body").item(0);
//     dom.onload();
// }
