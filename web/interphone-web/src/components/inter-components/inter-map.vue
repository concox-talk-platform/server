<template>
    <div id="map">
    </div>
</template>
<script>
    export default {
        name: 'intermap',

        data(){
            return{

            }
        },
        methods: {
                //  地图初始化
                createMap () {
                /* eslint-disable */
                // 创建Map实例
                var map = new BMap.Map("map")
                // 初始化地图,设置中心点坐标和地图级别
                map.centerAndZoom(new BMap.Point(114.07, 22.62), 11)
                //添加地图类型控件
                map.addControl(new BMap.MapTypeControl({
                    mapTypes:[BMAP_NORMAL_MAP, BMAP_HYBRID_MAP]
                }))
                // 设置地图显示的城市 此项是必须设置的
                map.setCurrentCity("深圳")
                //开启鼠标滚轮缩放
                map.enableScrollWheelZoom(true)
                /* eslint-enable */

                // 随机向地图添加25个标注
                var bounds = map.getBounds();
               var gps_data = JSON.parse(localStorage.getItem('device_list')).map(e =>{
                    if(e.hasOwnProperty('gps_data')){
                        return e.gps_data
                    }
                })
                var map_gps =gps_data.filter(Boolean)  
                for(var i =0; i<map_gps.length;i++){
                    var point = map_gps[i];
                    var marker = new BMap.Marker(point);
                    map.addOverlay(marker);

                }
            },            
        },
        mounted() {
            this.createMap ()
        },
    }

</script>
<style>
#map{
    height: 100%
}
</style>