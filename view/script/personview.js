  var call_ret;
  // 笔主详细信息
  function personal_description(user){
    //头像
    var a_img = document.getElementById("info-portrait");
    a_img.src=user["HeadImg"];
    a_img.style.backgroundSize = "cover";
    //昵称
    var a_nick = document.getElementsByClassName("info-name").item(0);
    a_nick.innerText = user["NickName"];
    //文章数量
    var p_note = document.getElementsByClassName("personal-note").item(0);
    p_note.innerText =parseInt(user["Note_sum"],10);
    //浏览 & 收藏量
    var p_view = document.getElementsByClassName("personal-view").item(0);
    p_view.innerText= parseInt(user["View_sum"],10);
    var p_save = document.getElementsByClassName("personal-save").item(0);
    p_save.innerText = parseInt(user["Col_sum"],10);

    var uid = localStorage.getItem("MyUserId");
     if (uid == user.Id){
        var manage = document.getElementById("manage_mynote");
        manage.addEventListener('click',function(e){
            document.location="mine.html?page=mynote";
        });
        manage.style.visibility="visible";
    }
}
//笔记类型过滤
var filtnote = "";
// 显示笔主的笔记
var notelist = new Array();
var catelist = new Array();
function personal_notelist(arr){
    /*
        arr[
            0:user      obj
            1:note      array
            2:category  array
        ]
    */
   console.log("arr:",arr);
   notelist = arr[1];
   catelist = arr[2];
   var attach = document.getElementById("personal-tb");
   attach.innerHTML="";
   //笔记为空时返回
   if(notelist == null){
    var notenum = document.getElementById("Cnt");
    notenum.innerText="0";
      return;
   }
   if(notelist){
       //显示列表框
       var list = document.getElementsByClassName("note-list").item(0);
       list.style.visibility="visible";
   }
    var root = document.getElementById("tp_minenotelist").contentWindow.document;

  
    var src = root.getElementById("personal-notelist");

    var is_person_center = false;
    try {
        personal_description(arr[0]);
    } catch (error) {
        is_person_center=true;
    }

    for(var i=0;i< notelist.length; i++){
        if(!catelist[i]){
            continue;
        }
        var CategoryName = catelist[i]["CategoryName"];  //科目分类
        var CategoryContain =catelist[i]["CategoryContain"];//科目细分
        //过滤显示科目，
        if(filtnote !="全部"  &&  CategoryName != filtnote){
            continue;
        }
        var dest = src.cloneNode(true);

         for(var j=0; j< dest.cells.length; j++){
            switch (dest.cells[j].id) {
                case "ArticleId":
                   var sb=dest.cells[j];
                   sb=sb.getElementsByTagName("p").item(0);
                    sb.innerText=notelist[i]["ArticleId"];
                    break;
                case "CategoryName":
                    var descripe = CategoryName+">"+CategoryContain;
                    if(notelist[i]["Remark"])
                        descripe += ">"+notelist[i]["Remark"];
                    if(!CategoryName)
                        descripe="暂无分类";
                  
                    dest.cells[j].innerText = descripe;
                    break;
                case "View_num":
                    dest.cells[j].innerText="浏览:"+notelist[i]["View_num"];
                    break;
                case "Collection":
                    dest.cells[j].innerText="收藏:"+notelist[i]["Collection"];
                    break;
                case "NickName":
                    dest.cells[j].innerText=notelist[i]["NickName"];
                    break;

            }
         }

        var tmp = dest.getElementsByClassName("view_del").item(0);
        var btn = dest.getElementsByClassName("viewbtn").item(0);
        var headimg = dest.getElementsByClassName("headimg").item(0);

        //当前页面不在个人中心
         if(!is_person_center){
            tmp.remove();
            btn.style.display="inline";
        }
        if(headimg){
            headimg.src="/"+notelist[i]["HeadImg"];
        }
        attach.appendChild(dest);
    }
    //回调
    if(call_ret){
        call_ret(arr);
    }
}

function req_note(notecate){
    filtnote=notecate;
    //tools.js
    var param = getQueryVariable("authorId");
    if(!param)
        param="";
    var route = '/getAuthorNote?authorId='+param; 
    request_route(route,personal_notelist);
}
function req_collectNote(notecate,called){
    filtnote=notecate;
    call_ret = called;
    var route = '/getCollectNotes?catename='+notecate;
    request_route(route,personal_notelist);
}