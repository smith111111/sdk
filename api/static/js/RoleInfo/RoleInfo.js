$(function () {
    $('#div').css("display","none")
    $('#setRoleActionDiv').css("display","none")
    $('#divUpdate').css("display","none")
    init()
    $('#btnRoleSearch').click(function () {
        var roleName=$('#txtSearchURoleName').val()
        var roleRemark=$('#txtSearchRoleRemark').val()
        var par ={
            "roleName":roleName,
            "roleRemark":roleRemark,
        }
        init(par)
    })
})

function init(par) {
    $('#rt').datagrid(
        {
            url: '/Admin/RoleInfo/GetRoleInfo',
            title: '角色数据表格',
            width: 700,
            height: 400,
            fitColumns: true, //列自适应
            nowrap: false,//设置为true，当数据长度超出列宽时将会自动截取
            idField: 'Id',//主键列的列明
            loadMsg: '正在加载角色的信息...',
            pagination: true,//是否有分页
            singleSelect: false,//是否单行选择
            pageSize:10,//页大小，一页多少条数据
            pageNumber: 1,//当前页，默认的
            pageList: [5,10, 20],
            queryParams:par,//往后台传递参数
            columns: [[
                //对应字段名字
                { field: 'ck', checkbox: true, align: 'left', width: 50 },//复选框设置
                { field: 'Id', title: '编号', width: 80 },
                { field: 'RoleName', title: '角色', width: 120 },
                { field: 'Remark', title: '备注', width: 120 },
                { field: 'AddDate', title: '时间', width: 80, align: 'right',
                    formatter: function (value, row, index) {
                        return value.split('T')[0]//对日期时间的处理 //看返回的时间进行割去
                    }
                }
            ]],
            toolbar: [{
                id: 'btnDelete',
                text: '删除',//显示的文本
                iconCls: 'icon-remove', //采用的样式
                handler: function () {	//当单击按钮时执行该方法
                    deleteRole();
                }
            },{
                id:'btnAdd',
                text:'添加',
                iconCls:'icon-add',
                handler:function () {
                    //展示添加用户表单
                    showAddRole()
                }
            },{
                id:'btnUpdate',
                text:'更新',
                iconCls:'icon-edit',
                handler:function () {
                    UpdateRole();//展示添加用户表单
                }
                },
                {
                    id:'btnRoleActions',
                    text:'角色分配权限',
                    iconCls:'icon-edit',
                    handler:function () {
                        AddRoleActions();//展示添加用户表单
                    },}
                    ],
        })
}

function AddRoleActions() {
//为角色分配权限
        var rows = $('#rt').datagrid('getSelections');
        if(rows.length!=1){
            $.messager.alert("提示","只能选择1个角色进行权限分配!!","error");
            return;
        }
        var roleId=rows[0].Id;
        $("#setRoleActionFrame").attr("src","/Admin/RoleInfo/ShowSetRoleAction?roleId="+roleId);
        $("#setRoleActionDiv").css("display","block");
        $('#setRoleActionDiv').dialog({
            title: '为角色分配权限信息',
            width: 300,
            height: 300,
            collapsible: true,
            maximizable: true,
            resizable: true,
            modal: true,
            buttons: [{
                text: 'Ok',
                iconCls: 'icon-ok',
                handler: function () {
                    //提交表单。
                    var childWindows=$("#setRoleActionFrame")[0].contentWindow;//获取子窗体的windows对象。
                    childWindows.formRoleActionSubmit()//提交表单

                }
            }, {
                text: 'Cancel',
                handler: function () {
                    $('#setRoleActionDiv').dialog('close');
                }
            }]
        });
}

function deleteRole() {
    var rows= $('#rt').datagrid('getSelections')
    if (rows.length==0){
        $.messager.alert("提示","请选择行删除","err")
        return
    }
    $.messager.confirm("提示","确认要删除1吗",function (r) {
        if(r){
            var str =""
            for (var i =0;i<rows.length;i++){
                str+=rows[i].Id+","//这里取得数据库返回的keyId
            }
            str=str.substring(0,str.length-1)
            var par={
                param:str
            }

            $.post("/Admin/RoleInfo/DeleteRole",par,function (data) {
                if (data.flag=="ok"){
                    $.messager.alert("提示","删除12成功","info")
                    $('#rt').datagrid('reload')
                    $('#rt').datagrid('clearSelections')
                }
            })

        }
    })



}

function showAddRole(){
    $('#addIframe').attr("src","/Admin/RoleInfo/GetRoleAddFormUrl")
    $('#div').css("display","block")
    $("#div").dialog({
        title: '新增角色信息',
        width: 300,
        height: 300,
        collapsible: true, //可折叠
        maximizable: true, //最大化
        resizable: true,//可缩放
        modal: true,//模态，表示只有将该窗口关闭才能修改页面中其它内容
        buttons: [{ //按钮组
            text: 'Ok',//按钮上的文字
            iconCls: 'icon-ok',
            handler: function () {
                var childWindow= $("#addIframe")[0].contentWindow
                childWindow.formSubmit()
            }
        }, {
            text: 'Cancel',
            handler: function () {
                $('#div').dialog('close');
            }
        }]
    });
}


function closeDialog(data) {
    if (data.flag == 'ok'){
        $('#div').dialog('close');
        $('#rt').datagrid('reload');
        $('#rt').datagrid('clearSelections')
    }
}
function closeUpdateDialog(data) {
    if (data.flag == 'ok'){
        $('#divUpdate').dialog('close');
        $('#rt').datagrid('reload');
        $('#rt').datagrid('clearSelections')
    }
}

function afterRoleAction(data) {
    if(data.flag=="yes"){
        $.messager.alert("提示","为角色分配权限成功!!","info")
        $('#setRoleActionDiv').dialog('close');
    }else{
        $.messager.alert("提示","为角色分配权限失败!!","error")
    }
    $('#setRoleActionDiv').dialog('close');
}

function UpdateRole() {
    //根据select的数据来比较
    var rows=$('#rt').datagrid('getSelections')
    if(rows.length!=1){
        $.messager.alert("提示","请选择一条数据修改","err")
        return
    }
    var par={
        Id:rows[0].Id
    }

    // //数据库查询数据比较更新的数据，不用rows,防止删除时还显示数据没加载
    $.post("/Admin/RoleInfo/SearRole",par,function(data){
        $('#updateIframe').attr("src","/Admin/RoleInfo/GetRoleFormUrl?id="+data.role.Id +
            "&roleName="+data.role.RoleName+"&remark="+data.role.Remark)
    })

    $('#divUpdate').css("display","block")
    $("#divUpdate").dialog({
        title: '编辑角色信息',
        width: 300,
        height: 300,
        collapsible: true, //可折叠
        maximizable: true, //最大化
        resizable: true,//可缩放
        modal: true,//模态，表示只有将该窗口关闭才能修改页面中其它内容
        buttons: [{ //按钮组
            text: 'Ok',//按钮上的文字
            iconCls: 'icon-ok',
            handler: function () {
                var childWindow= $("#updateIframe")[0].contentWindow//获取了ifrmae中的子窗体的window对象。
                childWindow.formUpdateSubmit()//调用子窗体中的方法。
            }
        }, {
            text: 'Cancel',
            handler: function () {
                $('#divUpdate').dialog('close');
            }
        }]
    });
}
