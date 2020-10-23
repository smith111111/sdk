$(function () {
    loadData()
    $('#useAction').css("display","none")
})

function loadData() {
    $('#tt').datagrid({
        url: '/Admin/ActionInfo/GetActionInfo',
        title: '权限数据表格',
        width: 970,
        height: 895,
        fitColumns: true, //列自适应
        nowrap: false,
        idField: 'Id',//主键列的列明
        loadMsg: '正在加载权限的信息...',
        pagination: true,//是否有分页
        singleSelect: false,//是否单行选择
        pageSize:10,//页大小，一页多少条数据
        pageNumber: 1,//当前页，默认的
        pageList: [2, 5, 10],
        queryParams: {},//往后台传递参数
        columns: [[//c.UserName, c.UserPass, c.Email, c.RegTime
            { field: 'ck', checkbox: true, align: 'left', width: 50 },
            //{ field: 'Id', title: '编号', width: 80 },
            { field: 'ActionInfoName', title: '权限名称', width: 120 },
            { field: 'HttpMethod', title: '请求方式', width: 120 },
            { field: 'Url', title: '请求地址', width: 120 },
            { field: 'Remark', title: '备注', width: 120 },
            { field: 'ActionTypeEnum', title: '权限类型', width: 120,
                formatter:function (value,row,index) {
                    return value=="1"?"菜单权限":"普通权限"
                }
            },
            { field: 'AddDate', title: '时间', width: 120, align: 'right',
                formatter: function (value, row, index) {
                    return value.split('T')[0]
                }
            }
        ]],
        toolbar: [{
            id: 'btnDelete',
            text: '删除',
            iconCls: 'icon-remove',
            handler: function () {
                showDeleteActionInfo()
            }
        },{
            id:'btnAdd',
            text:'添加',
            iconCls:'icon-add',
            handler:function () {
                showAddActionInfo();
            }
        }
        // ,
        //     {
        //     id:'btnEdit',
        //     text:'编辑',
        //     iconCls:'icon-edit',
        //     handler:function () {
        //
        //     }
        // }
        ],
    });
}

function showDeleteActionInfo() {
    //获取在表格中选中的行（getSelections：表示获取选中的行）
    var rows = $('#tt').datagrid('getSelections');
    if (!rows || rows.length == 0) {//判断是否选择了，如果没有选择长度为0
        $.messager.alert("提醒", "请选择要删除的记录!", "error");
        return;
    }
    $.messager.confirm("提示","确认要删除吗",function (r) {
        if (r){
            var strs = "";
            for (var i=0;i<rows.length;i++){
                strs += rows[i].Id+",";//Id大写
            }
            strs=strs.substring(0,strs.length-1);
            $.post('/Admin/ActionInfo/DeleteAction',{"ostrArr":strs},function (data) {
                if (data.flag=="ok"){
                    $.messager.alert("提示","删除菜单成功","info")
                    $('#tt').datagrid('clearSelections')
                    $('#tt').datagrid('reload')
                }else{
                    $.messager.alert("提示","删除菜单失败","info")
                }
            })
        }
    })
}

function showAddActionInfo() {
    $('#useActionIframe').attr("src","/Admin/ActionInfo/GetActionInfoAddFormUrl")
    $('#useAction').css("display","block")
    $("#useAction").dialog({
        title: '新增权限信息',
        width: 500,
        height: 500,
        collapsible: true, //可折叠
        maximizable: true, //最大化
        resizable: true,//可缩放
        modal: true,//模态，表示只有将该窗口关闭才能修改页面中其它内容
        buttons: [{ //按钮组
            text: 'Ok',//按钮上的文字
            iconCls: 'icon-ok',
            handler: function () {
                 var childWindow= $("#useActionIframe")[0].contentWindow
                 childWindow.formActionSubmit()
            }
        }, {
            text: 'Cancel',
            handler: function () {
                $('#useAction').dialog('close');
            }
        }]
    });
}

function closeActionDialog(data) {
    if(data.flag=="ok"){
        $.messager.alert("提示","添加成功","info");
        $('#useAction').dialog('close');
        $('#tt').datagrid('reload');//重新加载表格中的数据
    }else{
        $.messager.alert("提示","添加失败","error");
    }
}