      //identify 
    function initData(notecate) {
        req_note(notecate);
        categoryMyNote(notecate);
    }
    //我的笔记-上部基本信息
      function showData(e){
        var dom = document.getElementsByClassName("cagItem").item(0);
        var sp = dom.getElementsByTagName("p");
        for(var i=0;i<sp.length;i++){
            if(sp.item(i).id){
                if(e[sp.item(i).id] != "")
                    sp.item(i).innerText=e[sp.item(i).id];
            }
        }
    }
    //我的笔记-清除上部信息概要
    function clearData(area){
        var dom = document.getElementsByClassName("cagItem").item(0);
        var sp = dom.getElementsByTagName("p");
        for(var i=0;i<sp.length;i++){
            switch (sp.item(i).id) {
                case "Cnt":
                case "View_sum":
                case "Collect_sum":
                sp.item(i).innerText=0;
                break;
            } 
        }
        if(area == "all"){
            var sel_drop = document.getElementsByClassName("mycate").item(0);
            sel_drop.options.length=0;
        }
   
    }
    var sel_drop="";
    function categoryMyNote(notecate){
        request_route("/categoryMyNote",notecate_callback);
        function notecate_callback(e) {
             if(e == "未登录呢"){
              check_identify_jump("身份验证...","/login.html");
            }else{
                if(!e || e==""){
                  clearData("all");
                    return;
                }
                var Cnt=0;
                var View_sum=0;
                var Collect_sum=0;
                if(!sel_drop){
                    //添加select options
                    sel_drop = document.getElementsByClassName("mycate").item(0);
                    sel_drop.options.add(new Option("全部","全部"));
                    for(var i=0;i<e.length;i++){
                        if(!e[i])
                            continue;
                        //防止重复添加categoryName
                        var hasExist = false;
                        for(var j=0;j<sel_drop.options.length;j++){
                            if(sel_drop.options[j].value == e[i]["CategoryName"]){
                                hasExist = true;
                                break;
                            }
                        }
                        if(!hasExist){
                            var opt = new Option(e[i]["CategoryName"],e[i]["CategoryName"]);
                            sel_drop.options.add(opt);
                        }
                        
                    }
                }
                  for(var i=0;i<e.length;i++){
                     if(!e[i])
                        continue;
                    //过滤指定类型
                    if(e[i]["CategoryName"] == notecate){
                         showData(e[i]);
                    }
                     //if(notecate == "全部"){
                        Cnt += autoGetTypeData(e[i]["Cnt"]);
                        View_sum += autoGetTypeData(e[i]["View_sum"]);
                        Collect_sum += autoGetTypeData(e[i]["Collect_sum"]);
                    // }
                }
            
                //如果没有指定类型，数据计算后更新到e[0]
                if(notecate == "全部"){
                    var obj = e[0];
                    obj["CategoryName"]=notecate;
                    obj["Cnt"]=Cnt;
                    obj["View_sum"]=View_sum;
                    obj["Collect_sum"]=Collect_sum;
                    showData(obj);
                }
                //选中不存在的分类后清除信息
                var clear = true
                for(var j=0;j<e.length;j++){
                    if(e[j].CategoryName == notecate){
                        clear=false
                        break;
                    }
                }
                if(clear){
                    clearData();
                }
                
            }
        }

    }

              //删除笔记成功回调
              function delt_callback(noteid,deltdom,data){
                console.log("deal_callback:"+data);
                if(data == "success" && deltdom){
                    deltdom.remove();
                    var sel = document.getElementsByClassName("mycate").item(0);
                    //重新请求数据
                    initData(sel[sel.selectedIndex].value);
                }else{
                    show_warn("失败",data);
                }
            }
            //删除 or 查看
             function deal(dom,ope) {
             if(dom.id == "ArticleId"){
                 var noteid = dom.getElementsByTagName("p").item(0).innerText;
                 if(ope == "delt_note" || ope=="colcel_note"){
                     //弹出确认框
                     show_warn("警告","是否确认删除?", function(){
                     //已确认，发送请求
                     send_data(ope, function(e){
                       delt_callback(noteid,dom.parentElement, e);
                     }, "noteId=" + noteid);
                     });
                 }else{
                    var catename = dom.parentElement.getElementsByClassName("catename").item(0).innerText;
                    encode_catename = encodeURI(encodeURI(catename));
                      document.location="/notedetail.html?noteid="+noteid+
                     "&ope="+ope+"&catename="+encode_catename;
                 }
             return;
             }else{
                 deal(dom.parentElement,ope);
             }
   
         }