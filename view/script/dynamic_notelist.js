//传递一个请求路由，得到服务器返回的json
function create_note_item(){
    var attach=document.getElementById("note_list_item");  
    var impnoteItem = document.getElementById("noteItem").contentWindow.document;
    var src = impnoteItem.getElementsByClassName("item").item(0);
    var item = src.cloneNode(true);
    var sle="visibility: visible;";
    item.style=sle;
    attach.appendChild(item);
    return item;
}
/*
    arr(
        map[author,note]
    )
*/
 
function setData(arr){
    var arr_len = arr.length;
    if(arr_len < 1)
        return;
    //当前列表总长度 i -> all_len
    for(var i=0; i<arr_len; i++){
        var author = arr[i]["author"];
        var note = arr[i]["note"];
        if(!author || note.length < 1){
            return;
        }
        var cl1 = create_note_item();
        var el= cl1.getElementsByTagName("label");
        var divele=cl1.getElementsByClassName("head").item(0);
        //console.log("arr:",arr[ix]);
        var styles = "background:url("+author.HeadImg+") no-repeat; \
        background-size:cover;cursor: pointer;";
        divele.style=styles;    
        //对每个Item子项元素赋值
          for (j=0;j<el.length;j++){
             switch (el.item(j).id) {
                case "authorId":
                    el.item(j).innerText=author.Id; 
                    break;
                case "ArticleId":
                    el.item(j).innerText=note.ArticleId;
                    break;
                case "project":
                    el.item(j).innerText=note.CategoryName+"笔记";
                    break;
                case "authorName":
                    el.item(j).innerText=author.NickName;
                     break;
                case "description":
                    el.item(j).innerText="简介:"+author.Description;
                    break;
                case "catedetail":
                    var cate = "分类:"+note.CategoryContain;
                    if(note.Remark){
                        cate +=">"+note.Remark
                    }
                    el.item(j).innerText=cate;
                    break;  
                case "view_num":
                    if(note.View_num == "")
                        note.View_num=0;
                    el.item(j).innerText="浏览:"+note.View_num;
                    break;
             }
        }
        //详细信息按钮
        var hebtn = cl1.getElementsByClassName("head").item(0);
        hebtn.addEventListener("click",viewDetail);
        var debtn = cl1.getElementsByClassName("descibe_btn").item(0);
        debtn.addEventListener("click",viewNote);

    }

}
    //查看笔记详情
    function viewNote(){
        var dom = this.parentNode;
        var noteid = dom.getElementsByClassName("ArticleId").item(0).innerHTML;
        var catename = dom.getElementsByClassName("catedetail").item(0).innerHTML;
        var url="/notedetail.html?ArticleId="+noteid+
        "&ope=view_note"+"&catename="+catename;
        window.open(url,"_blank");
    }
    function viewDetail(){
        var authorId = this.getElementsByClassName("authorId").item(0).innerHTML;
        window.open("personview.html?authorId="+authorId,"_blank");
     }

function menu_item2_mouseleave(){
    var menu=document.getElementById("expandMenu");
    var extenditem=document.getElementById("expandItem");
    if(isShouldHideMenu("expandMenu")){
         menu.innerHTML="";
    }else{
        extenditem.innerHTML="";
    }
       
}
function menu_item_mouseleave(){
    /*
    隐藏扩展菜单条件
    离开扩展菜单后是否进入导航菜单，如果鼠标划过扩展菜单的其它项,
    它仍然触发此事件，此时不应该隐藏扩展菜单
    */
    if(isShouldHideMenu("brand") && isShouldHideMenu("expandMenu") && isShouldHideMenu("expandItem")){
        var menu=document.getElementById("expandMenu");
        menu.innerHTML="";
    }
}
function menu_mouseleave(){
    //离开导航菜单后鼠标是否进入扩展菜单中
    if(isShouldHideMenu("expandMenu")){
        var menu=document.getElementById("expandMenu");
        menu.innerHTML="";
    }
 }
 function isShouldHideMenu(eleid){
    var x= window.event.clientX;
    var y=window.event.clientY;
    var ele=document.getElementById(eleid);
    var eleX=Math.floor(ele.getBoundingClientRect().left);
    var eleY=Math.floor(ele.getBoundingClientRect().top);
    var eleW=ele.getBoundingClientRect().width;
    var eleH=ele.getBoundingClientRect().height;
    if(x<eleX || x>eleX+eleW ||y< eleY || y>eleY+eleH){
       return true;
    }
    return false;
 }
 
function menu_item_enter(){
    var menu=document.getElementById("expandMenu");
    var itemContain=document.getElementById("expandItem");
    if(!itemContain){
        itemContain=document.createElement("div");
        itemContain.id="expandItem";
        menu.appendChild(itemContain);
    }
    itemContain.innerHTML="";
    itemContain.addEventListener("mouseleave",menu_item2_mouseleave);
    //当前学科
    var proj=new Category();
    proj.contain = new Array();
    for(var i = 0; i < cate.length; i++){
        //去空格
        var v1 = cate[i].name.replace(/^\s+|\s+$/g, '');
        var v2 = this.innerText.replace(/^\s+|\s+$/g, '');
        if(v1 == v2){
            proj=cate[i];
            break;
        }
    }
    //学科->所有年级
    var rankcnt = proj.contain.length;
    for(var i=0; i < rankcnt; i++){
        var item=document.createElement("div");
        //弹出菜单上边距依附当前对象
        item.id="detailItem";
        var ofy = this.offsetTop;
        itemContain.style.top=ofy;
        item.innerText=proj.contain[i];
        itemContain.appendChild(item);
    }
    if(rankcnt == 0){
        alert("不支持edge");
    }
  
   
}
//中部，创建导航菜单
function createMenu(_this,func){
 
    var menu=document.getElementById("expandMenu");
    menu.innerHTML="";
    var arrItem=new Array();
    if(_this.innerText == "学科分类"){
           arrItem=cate;
    }
    var len=arrItem.length;
    for(i=0;i<len;i++){
    var menu_item=document.createElement("ul");
    menu_item.id="menuItem";
    menu_item.innerText=cate[i].name;
   menu_item.addEventListener('mouseenter',menu_item_enter);
   menu_item.addEventListener('mouseleave',menu_item_mouseleave)
    var menu_item_line=document.createElement("label");
    menu_item_line.className="separate";
    menu_item.appendChild(menu_item_line);
    menu.appendChild(menu_item);
    }
}
