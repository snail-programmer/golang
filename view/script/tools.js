function converToNumber(params) {
    var res;
    if(params.indexOf(".",0)<0){
        res = parseInt(params,10);
    }else{
        res = parseFloat(params);
    }
    return res;
}
function autoGetTypeData(data) {
    var re = /^[0-9]+.?[0-9]*$/; 
    var resp=data;
    if(re.test(data)){
        resp = converToNumber(data);
    }
    return resp;

}
function getCurTime(){
    var myDate = new Date();
    var datetime= myDate.getFullYear()+"/"+
    myDate.getMonth()+"/"+myDate.getDate()+" "+
    myDate.getHours()+":"+ myDate.getMinutes()+":"+ myDate.getSeconds();
    return datetime;
}

//xml deal
function parseDom(str) {
    var div = document.createElement("div");
    div.innerHTML = str;
    return div;
}
function getxml(route,searchId) {
    var xml = new XMLHttpRequest();
    xml.open("GET",route,false);
    xml.send();
        if(xml.readyState == 4 && xml.status==200){
        var xmlstr= xml.response;
        var xmlparse = new DOMParser();
        var xmldoc = xmlparse.parseFromString(xmlstr,"text/xml");
        var dom = xmldoc.getElementById(searchId);
        var cont = dom.getElementsByTagName("content").item(0);
         return cont.innerHTML;
     }  
}
function getQueryVariable(variable)
{
       var query = window.location.search.substring(1);
       var vars = query.split("&");
       for (var i=0;i<vars.length;i++) {
               var pair = vars[i].split("=");
               if(pair[0] == variable){return pair[1];}
       }
       return(false);
}