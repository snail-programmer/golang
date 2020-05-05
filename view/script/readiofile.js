function readiofile(domfile,callback){
    domfile.onchange=function(){
        var read = new FileReader();
        read.readAsArrayBuffer(domfile.files[0]);
        read.onload=function(){
             if(callback){
                var e = this.result;
                var dgk = new TextDecoder('gbk').decode(e);
                var dut = new TextDecoder().decode(e);
                var lk=dgk.length;
                var lt=dut.length;
                var real = lk < lt?dgk:dut;
                 callback(real);
             }
        }
    }
}