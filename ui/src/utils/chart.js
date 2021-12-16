export const formatTooltip = (params,unit)=>{
    const val = formatUnit(params.data.value[1],unit)
    const date =params.data.value[0]
    return `<div style="margin: 0px 0 0;line-height:1;">
                <div style="font-size:14px;color:`+params.color+`;font-weight:900;line-height:1;">`+params.seriesName+`</div>
                <div style="margin: 10px 0 0;line-height:1;">
                    <div style="margin: 0px 0 0;line-height:1;">
                        <span style="display:inline-block;margin-right:4px;border-radius:10px;width:10px;height:10px;background-color:`+params.color+`;"></span>
                         <span style="float:right;margin-left:10px;font-size:14px;color:#666;font-weight:900">`+date+`</span>
                         <span style="float:right;margin-left:10px;font-size:14px;color:#666;font-weight:900">`+val+`</span>
                        <div style="clear:both"></div>
                    </div>
                    <div style="clear:both"></div>
                </div>
                <div style="clear:both"></div>
             </div>`
}

export const formatUnit = (value,unit)=>{
    if (unit === "bytes"){
        if (null == value || value == '') {
            return "0 Bytes";
        }
        var unitArr = new Array("Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB");
        var index = 0;
        var srcsize = parseFloat(value);
        index = Math.floor(Math.log(srcsize) / Math.log(1024));
        var size = srcsize / Math.pow(1024, index);
        size = size.toFixed(2);//保留的小数位数
        return size + unitArr[index];
    }else if (unit === "nanoseconds"){
        if (null == value || value == '') {
            return "0";
        }
        var unitArr = new Array("ns", "us", "ms", "s");
        var index = 0;
        var srcsize = parseFloat(value);
        index = Math.floor(Math.log(srcsize) / Math.log(1000));
        if (index > 3){
            index = 3
        }
        var size = srcsize / Math.pow(1000, index);
        size = size.toFixed(2);//保留的小数位数
        return size + unitArr[index];
    }
    else if (unit === "count"){
        return value
    }
    return value+unit
}
