var globaldom;
var dom;
function initMask(){
    var link = document.createElement("link");
    link.rel = "stylesheet";
    link.type = "text/css";
    link.href = "/css/mask.css";
    globaldom = window.top.document;
    var head = globaldom.getElementsByTagName("head")[0];
    head.appendChild(link);
    dom = globaldom.getElementsByTagName("body").item(0);
 }
function show_mask(contStr){
    initMask();
    var mask = document.createElement("div");
    var content= document.createElement("div");
    var img = document.createElement("img");
    var info = document.createElement("p");
    mask.className="_screen_mask";
    content.className="_mask_content";
    img.className="loadimg";
    img.src="../image/largeloading.gif";
    info.innerText="加载中...";
    if(contStr){
        info.innerText = contStr;
    }
    if(dom && mask){
        content.appendChild(img);
        content.appendChild(info);
        dom.appendChild(mask);
        dom.appendChild(content);
    }
 }
 
function show_maskTime(contStr,time,callback) {
    show_mask(contStr);
    setTimeout(() => {
        hide_mask();
        if(callback)
            callback();
    }, time);
}
function hide_mask(){
    
   var mask = globaldom.getElementsByClassName("_screen_mask").item(0);
    var content= globaldom.getElementsByClassName("_mask_content").item(0);
    if(mask && content){
        mask.remove();
        content.remove();
    }
}