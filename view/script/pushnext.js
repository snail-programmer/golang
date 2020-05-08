var topdcm=window.top.document;
var pn_cate;
//初始化第一个分类
function init_categoryName() {
    var cate = pn_cate;
   var catename = topdcm.getElementsByClassName("CategoryName").item(0);
   catename.options.length=0;
   var vkey = "";
   for(var key in cate){
    var option = new Option(key,key);
    catename.options.add(option);
    if(vkey == "")
        vkey = key;
   }

   init_categoryContain();
}
//初始化第二个分类
function init_categoryContain() {
    var cate = pn_cate;
   var catename = topdcm.getElementsByClassName("CategoryName").item(0);
   var cate_sel = catename[catename.selectedIndex].value;
   var contain;
   for(var key in cate){
       if(key == cate_sel){
           contain = cate[key];
           break;
       }
   }
   var catecontain =  topdcm.getElementsByClassName("CategoryContain").item(0);
   catecontain.options.length=0;
   for(var i=0; i<contain.length;i++){
       catecontain.options.add(new Option(contain[i],contain[i]));
   }
}
function show_pushnext() {
   var showcate = window.top.document.getElementById("showcate").contentWindow.document;
   var back = showcate.getElementById("pushnext");
   var newback=back.cloneNode(true);
   var btn = newback.getElementsByClassName("pushnextbtn").item(0);
   btn.addEventListener("click",pushnext_btn);
   var second = newback.getElementsByClassName("CategoryName").item(0);
   second.addEventListener("change",init_categoryContain);

   var attach=window.top.document.getElementsByTagName("body").item(0);
   attach.appendChild(newback);

   //请求分类 ->category ->cate array
   request_route('/categorylist',function(e){
        pn_cate=e;
        init_categoryName();
   });

 
}
function hide_pushnext() {
   var back = window.top.document.getElementById("pushnext");
    back.remove();
}

          //成功回调
          function pushNote_callback(data){
            hide_mask();
            console.log("cal:",data);
            show_warn("成功","笔记发布完成!");
        }
        function pushnext_btn() {
            
            show_mask("正在发布...");
            var articleId=document.getElementById("ArticleId").innerText;
            var publishTime= getCurTime();
            var categoryName= topdcm.getElementsByClassName("CategoryName").item(0);
            var categoryContain=topdcm.getElementsByClassName("CategoryContain").item(0);
            var remark = topdcm.getElementsByClassName("remark").item(0);
            categoryName = categoryName[categoryName.selectedIndex].value;
            categoryContain= categoryContain[categoryContain.selectedIndex].value;
            remark = remark.value;
            var data="createTime="+publishTime+"&article="+$("#EditNote").val()+
            "&categoryName="+categoryName+"&categoryContain="+categoryContain+"&remark="+remark;
            if(articleId){
                data+="&articleId="+articleId;
            }
            send_data('/xhedit_saveNote',pushNote_callback,data);
           hide_pushnext();

        }