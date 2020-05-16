function create_ele(divname,propname){
    var ele = document.createElement(divname);
    if(propname){
        var idcs = propname.split("=");
        if(idcs[0]=="class"){
            ele.className=idcs[1];
        }else if(idcs[0]=="id"){
            ele.id=idcs[1];
        }else{
            ele.style=propname;
        }
    }
    return ele;
}
function init_comment(){
    var content = document.getElementsByClassName("comment-content").item(0);
    content.innerHTML="";
    request_route("/get_comment?ArticleId="+aid,function(e){
        console.log("comment:",e);
        create_comment(e);
        //位置跳转
        var jump = getQueryVariable("jump");
        if(jump != "false"){
             location.href="#"+jump;
        }
    });
}
function create_comment(arr){
    var refary = new Array();
    var content = document.getElementsByClassName("comment-content").item(0);
    var myid = localStorage.getItem("MyUserId");
    for(var i=0;i<arr.length; i++){
        if(!arr[i])
            break
        var item = create_ele("div","class=comment-item");
        var head = create_ele("div","class=comment-head");
        var id = create_ele("label","display:none;");
        var img =create_ele("img","class=headimg");
        var nick = create_ele("label","class=nickname");
        var cope = create_ele("button","class=cmt-ope");
        var time = create_ele("label","class=comment-time");
        var comment = create_ele("p","class=comment-info");
        img.src=arr[i].HeadImg;
        id.className="myid";
        item.id=arr[i].Id;
        id.innerText=arr[i].MyId;
        nick.innerText=arr[i].NickName;
        time.innerText=arr[i].Cmt_time;
        comment.innerHTML=arr[i].Comment;
        head.appendChild(id);
        head.appendChild(img);
        head.appendChild(nick);
        head.appendChild(cope);
        head.appendChild(time);
        item.append(head);
        if(arr[i].ReplyId){
            var refer = create_ele("div","class=refer-info");
            refary[i]=refer;
                request_route("/get_comment?id="+arr[i].ReplyId,function(e){
                    if(!e[0]){
                        referData =e.msg;
                    }else{
                        referData = "@<a target='_blank' href='/personview.html?authorId="+e[0].MyId+"'>"+e[0].NickName+"</a></br>"+e[0].Comment;
                    } 
                    for(var ri =0;ri<refary.length;ri++){
                        if(refary[ri]){
                            refary[ri].innerHTML=referData;
                            refary[ri]=null;
                            break;
                        }
                    } 
                });
            item.appendChild(refer);
        }
        comment.innerHTML=arr[i].Comment;
        item.append(comment);
        content.appendChild(item);
        if(myid != arr[i].MyId){
            cope.innerText="回复";
            cope.addEventListener("click",function(){
            var doc=this.parentElement.parentElement;
            var identify = doc.id;
            var id = doc.getElementsByClassName("myid").item(0).innerText;
            var nick = doc.getElementsByClassName("nickname").item(0).innerText;
            var cmt = doc.getElementsByClassName("comment-info").item(0).innerText;
            var warn = new CWarnInfo("input");
            warn.show_warn("回复@"+nick,"neir",function(e){
                e+="&replyId="+identify;
                send_comment(e);
            });
        });
        }else{
            cope.innerText="删除";
            cope.addEventListener('click',function(){
                var cmt_time=this.parentElement.getElementsByClassName("comment-time").item(0);
                var url="/del_comment?cmt_time="+cmt_time.innerText;
                show_warn("提示","删除评论?",function(){
                    request_route(url,function(e){
                    show_warn("提示",e.msg);
                    init_comment();
                });
                });
               
            });
        }
 
    }
}