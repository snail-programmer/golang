var globaldom;
var dom;
var global_mask;
var mask_content;
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
    mask.style.top=document.body.scrollTop+document.documentElement.scrollTop+"px";
    var scrollTop=document.body.scrollTop+document.documentElement.scrollTop+140;
    content.style.top=scrollTop+"px";
    global_mask=mask;
    mask_content=content;
    if(contStr){
        info.innerText = contStr;
    }
    if(dom && mask){
        content.appendChild(img);
        content.appendChild(info);
        dom.appendChild(mask);
        dom.appendChild(content);
        window.top.document.addEventListener("scroll", scroll_listen,true);
    }
 
    function scroll_listen(){
        var scrollTop=document.body.scrollTop+document.documentElement.scrollTop;
        global_mask.style.top=scrollTop+"px";
        mask_content.style.top=scrollTop+140+"px";
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