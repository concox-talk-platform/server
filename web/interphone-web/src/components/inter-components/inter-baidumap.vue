<template>
    <div id="map" >
    </div>
</template>
<script>
    export default {
        name: 'intermap',

        // props:['postSOS'],
    
        data(){
            return{
                mapcontrol :''

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
                this.mapcontrol =map
            //   var bounds = map.getBounds();
            //    var gps_data = JSON.parse(localStorage.getItem('device_list')).map(e =>{
            //         if(e.hasOwnProperty('gps_data')){
            //             return e.gps_data
            //         }
            //     })
            //     var map_gps =gps_data.filter(Boolean)  
            var map_gps =[{lng: 113.98336974606138, lat: 22.695076044933757},
                 {lng: 113.94189585189847, lat: 22.510310178430622},{lng: 113.32123876460987, lat: 22.808882485358566},
                ];
                for(var i =0; i<map_gps.length;i++){
                    var point = map_gps[i];
                    var marker = new BMap.Marker(point);
                    map.addOverlay(marker);

                }
                // var b =new BMap.Point( 114.1153786184, 22.4644910453);
                // var a =new BMap.Marker(b)
                // map.addOverlay(a)
                // a.setAnimation(BMAP_ANIMATION_BOUNCE);
                // window.console.log('sos`````````````````')

            },     
            sos_point(e){ 
                var b =new BMap.Point(e.lng, e.lat);
                var a =new BMap.Marker(b)
                this.mapcontrol.addOverlay(a)
                a.setAnimation(BMAP_ANIMATION_BOUNCE);
                window.console.log('sos`````````````````')
            }       
        },
        mounted() {
            this.createMap ();
            // window.console.log(this.postSOS)
            // this.sos_point()
        },

    }

</script>
<style>
#map{
    height: 100%
}
</style>