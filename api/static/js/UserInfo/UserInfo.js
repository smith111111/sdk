$(function () {
    $("#addDiv").css("display","none")
    $("#setUserActionDiv").css("display","none")
    $('#btnSearch').click(function () {
        var txtSname=$('#txtSearchUName').val()
        var txtSremark= $('#txtSearchRemark').val()
        var parame ={
            "txtSname":txtSname,
            "txtSremark":txtSremark
        }
        loadData(parame)
    });
    loadData()
})
function loadData(parame){
    $('#tt').datagrid({
        url: '/Admin/UserInfo/GetUserInfo',
        title: '用户数据表格',
        width: 957,
        height: 871,
        fitColumns: true, //列自适应
        nowrap: false,//设置为true，当数据长度超出列宽时将会自动截取
        idField: 'Id',//主键列的列明
        loadMsg: '正在加载用户的信息...',
        pagination: true,//是否有分页
        singleSelect: false,//是否单行选择
        pageSize:10,//页大小，一页多少条数据
        pageNumber: 1,//当前页，默认的
        pageList: [2,10, 20],
        queryParams: parame,//往后台传递参数
        columns: [[
            //对应字段名字
            { field: 'ck', checkbox: true, align: 'left', width: 50 },//复选框设置
            { field: 'Id', title: '编号', width: 80 },
            { field: 'UserName', title: '姓名', width: 120 },
            { field: 'UserPwd', title: '密码', width: 120 },
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
                deleteUser();
            }
        },{
            id:'btnAdd',
            text:'添加',
            iconCls:'icon-add',
            handler:function () {
                showAddUser();//展示添加用户表单
            }
        },{
            id:'btnUpdate',
            text:'更新',
            iconCls:'icon-edit',
            handler:function () {
                UpdateUser();//展示添加用户表单
            }},
            {
                id:'btnUserAndRole',
                text:'用户分配角色',
                iconCls:'icon-edit',
                handler:function () {
                    UserAndRole();//展示添加用户表单
                },
            }, {
                id:'btnUserAndActions',
                text:'用户分配权限',
                iconCls:'icon-edit',
                handler:function () {
                    SetUserAction();//展示添加用户表单
                },
            }
            ],
    })
}

function SetUserAction(){
    var rows = $('#tt').datagrid('getSelections');
    if(rows.length!=1){
        $.messager.alert("提示","请选择要分配权限的用户!!","error");
        return;
    }
    $("#setUserActionFrame").attr("src","/Admin/UserInfo/ShowUserAction?userId="+rows[0].Id);
    $("#setUserActionDiv").css("display","block");
    $('#setUserActionDiv').dialog({
        title: '为用户分配角色',
        width: 850,
        height: 500,
        collapsible: true,
        maximizable: true,
        resizable: true,
        modal: true,
        buttons: [{
            text: 'Ok',
            iconCls: 'icon-ok',
            handler: function () {
                $('#setUserActionDiv').dialog('close');
            }
        }, {
            text: 'Cancel',
            handler: function () {
                $('#setUserActionDiv').dialog('close');
            }
        }]
    });
}


//删除用户

function deleteUser() {
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
            $.post('/Admin/UserInfo/DeleteUser',{"ostrArr":strs},function (data) {
                if (data.flag=="ok"){
                    $.messager.alert("提示","用户提交成功","info")
                    $('#tt').datagrid('clearSelections')
                    $('#tt').datagrid('reload')
                }else{
                    $.messager.alert("提示","用户提交失败","info")
                }



            })
        }
    })
}

function UpdateUser() {
    var rows = $('#tt').datagrid('getSelections');
    // if (!rows || rows.length == 0) {//判断是否选择了，如果没有选择长度为0
    //     $.messager.alert("提醒", "请选择要修改的商品!", "error");
    //     return;
    // }
    // if (rows.length > 1) {//判断是否选择了，如果没有选择长度为0
    //     $.messager.alert("提醒", "请选择一个商品修改!", "error");
    //     return;
    // }
    if (!rows || rows.length != 1) {//判断是否选择了，如果没有选择长度为0
        $.messager.alert("提醒", "请选择一条要修改的用户!", "error");
        return;
    }

    showAddUser(rows[0].Id)

}
function showAddUser(id) {

    $.post("/Admin/UserInfo/FindUserById",{'Id':id},function (data) {
        //判断服务端返回的数据。
        $('#addForm').form('load', data);

    })


    $('#addDiv').css("display","block");//显示
    $('#addDiv').dialog({
        title:(id? '更新':'添加')+'用户信息',
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
                //1:获取表单中的数据
                //2:发送服务端。
                // (id ? AddUserData():UpdateData());
                AddUserData()
            }
        }, {
            text: 'Cancel',
            handler: function () {
                $('#addDiv').dialog('close');
            }
        }]
    });
}

//获取表单中数据，完成添加
function AddUserData() {
    //1：获取表单中的数据
    var pars= $("#addForm").serializeArray()
    var user=pars[1]["value"];
    var pwd=pars[2]["value"];
    if(user==""){
        $.messager.alert("提示","用户名称不能为空!","info");
        return
    }
    if(pwd==""){
        $.messager.alert("提示","用户密码不能为空!","info");
        return
    }


    //2:发送数据/
    $.post("/Admin/UserInfo/AddUser",pars,function (data) {
        //判断服务端返回的数据。
        if (data.flag=="ok"){
            if(data.id!="") {
                $.messager.alert("提示","用户信息修改成功!!","info");
            }else{
                $.messager.alert("提示","用户信息添加成功!!","info");
            }
            $('#addDiv').dialog('close');
            $("#addForm input").val("")//将form中的所有的input标签赋空值。
            $('#tt').datagrid('reload')
            //loadData()
        } else{
            if(data.id!="") {
                $.messager.alert("提示","用户信息修改失败!!","info");
            }else{
                if (data.flag=="新增账号已经存在"){
                    $.messager.alert("提示",data.flag,"info");
                    return
                }
                $.messager.alert("提示","用户信息添加失败!!","info");
            }
        }
    })
}

function UserAndRole() {
    //根据select的数据来比较
    var rows=$('#tt').datagrid('getSelections')
    if(rows.length!=1){
        $.messager.alert("提示","请选择一条数据进行用户分配","err")
        return
    }
    // //数据库查询数据比较更新的数据，不用rows,防止删除时还显示数据没加载 这里只传个就可以
        $('#userRoleIframe').attr("src","/Admin/UserInfo/GetUserRoleFormUrl?id="+rows[0].Id)

    $('#userRole').css("display","block")
    $("#userRole").dialog({
        title: '分配角色信息',
        width: 400,
        height: 300,
        collapsible: true, //可折叠
        maximizable: true, //最大化
        resizable: true,//可缩放
        modal: true,//模态，表示只有将该窗口关闭才能修改页面中其它内容
        buttons: [{ //按钮组
            text: 'Ok',//按钮上的文字
            iconCls: 'icon-ok',
            handler: function () {
                var childWindow= $("#userRoleIframe")[0].contentWindow//获取了ifrmae中的子窗体的window对象。
                childWindow.formRoleSubmit()//调用子窗体中的方法。
            }
        }, {
            text: 'Cancel',
            handler: function () {
                $('#userRole').dialog('close');
            }
        }]
    });
}

function closeRoleDialog(data) {
    if (data.flag == 'yes'){
        $('#userRole').dialog('close');
    }

}