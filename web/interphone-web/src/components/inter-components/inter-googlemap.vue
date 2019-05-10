<template>
    <div id="googleMap"></div>
</template>

<script>
    import google from 'google'
    export default {
        name: 'intergooglemap',

        // props:['postSOS'],
    
        data(){
            return{
                google_map:''
            }
        },
        methods: {
            createmap(){
                var mapProp = {
                    // center:new google.maps.LatLng(51.508742,-0.120850),
                    // center:new google.maps.LatLng(39.5427,116.2317),
                    center:new google.maps.LatLng(22.62, 114.07),
                    zoom:11,
                    mapTypeId:google.maps.MapTypeId.ROADMAP
                };
                var map=new google.maps.Map(document.getElementById("googleMap"), mapProp);
                this.google_map =map;
                // var gps_data = JSON.parse(localStorage.getItem('device_list')).map(e =>{
                //     if(e.hasOwnProperty('gps_data')){
                //         return e.gps_data
                //     }
                // })
                // var map_gps =gps_data.filter(Boolean)  
                // {lng: 22.4644910453, lat: 114.1153786184},
                var map_gps =[ {lng: 113.98336974606138, lat: 22.695076044933757},
                {lng: 113.94189585189847, lat: 22.510310178430622},{lng: 113.32123876460987, lat: 22.808882485358566},
                ];
                var gps_place =[]
                for(var i = 0; i<map_gps.length;i++){
                    gps_place.push(new google.maps.LatLng(map_gps[i].lat,map_gps[i].lng));
                }
                window.console.log(gps_place);
                for (var k = 0; k < gps_place.length; k++) {
                    var marker = new google.maps.Marker({
                        position: gps_place[k],
                        map: map,
                        title: 'Place number ' + k
                    });
                  }
                },
                sos_google(e){
                  window.console.log(e)
                  var google_gps=new google.maps.LatLng(e.lat,e.lng)
                    var google_sos=new google.maps.Marker({
                        position: google_gps,
                        map: this.google_map,
                        animation:google.maps.Animation.BOUNCE,

                    });
                }
        },
        mounted() {
            this.createmap()
        },

    }

</script>
<style>
#googleMap{
    height: 100%
}
</style>