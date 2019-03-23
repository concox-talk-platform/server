<template>
    <div>
        <el-table ref="multipleTable" :data="tableData.slice((currentPage-1)*pagesize,currentPage*pagesize)" tooltip-effect="dark" style="width: 100%"  @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" ></el-table-column>
            <el-table-column  type="index" width="50" label="序号"></el-table-column>
             <el-table-column prop="imei" label="IMEI" width="240"> </el-table-column>
             <el-table-column prop="Equipment" label="设备" width="240"> </el-table-column>
             <el-table-column prop="staff" label="人员" width="240"> </el-table-column>
             <el-table-column prop="group" label="对讲组" width="240"> </el-table-column>
             <el-table-column prop="EquipmentState" label="设备状态" width="120"> </el-table-column>
             <el-table-column label="操作">
                   <template slot-scope="scope">
                             <el-button size="mini" @click="handleEdit(scope.$index, scope.row)">编辑</el-button>
                             <el-button size="mini" @click="handleUnbind(scope.$index, scope.row)">解绑</el-button>
                             <el-button size="mini" @click="handleAddgroup(scope.$index, scope.row)">添加到组</el-button>
                             <el-button size="mini" @click="handleMore(scope.$index, scope.row)">更多</el-button>
                             <el-button size="mini" type="danger" @click="handleDelete(scope.$index, scope.row)">删除</el-button>
                    </template>
             </el-table-column>
        </el-table>
        <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange" :current-page="currentPage" :page-sizes="[10, 20, 30, 40]"
         :page-size="10" layout="total, sizes, prev, pager, next, jumper" :total="400">
        </el-pagination>
    </div>
</template>
<script>
export default {
 data() {
      return {
        tableData: [{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },{
         imei: '21424152312321312',
          Equipment: '王小虎',
          staff: '测试组一',
          group: '对讲组10',
          EquipmentState: '离线',
    
        },],
        multipleSelection: [],
        currentPage: 1,
        pagesize:10,
        
      }
    },
    methods:{
        handleSelectionChange(val) {
          this.multipleSelection = val;
        },
        handleEdit(index, row) {
            window.console.log(index, row);
            // window.console.log(this.multipleSelection);
            
        },
        handleUnbind(index, row) {
            window.console.log(index, row);
        },
        handleAddgroup(index, row) {
            window.console.log(index, row);
        },
        handleMore(index, row) {
            window.console.log(index, row);
        },                
        handleDelete(index, row) {
            window.console.log(index, row);
        },
        handleSizeChange(val) {
            window.console.log(`每页 ${val} 条`);
            // this.currentPage = currentPage;
             this.pagesize = val;
        },
        handleCurrentChange(currentPage) {
        //    window.console.log(`当前页: ${val}`);
           this.currentPage =currentPage
        }
    }
}
</script>
<style>

</style>
