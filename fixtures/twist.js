
n() {

        // ------------------------------------
        // Bootstrap tooltip trigger
        // ------------------------------------
        $("a[data-toggle=tooltip]").tooltip();


        // ------------------------------------
        // Twist MapManager
        // ------------------------------------

        var mapCenter = new google.maps.LatLng(45.4642700, 9.18951);


 
var twistMap = new TwistMapManager('map', mapCenter, {
            blueCarPoints: [ 
new google.maps.LatLng(45.453168468953,9.1660584806154),
new google.maps.LatLng(45.453297990697,9.1734399916735),
new google.maps.LatLng(45.434840355723,9.1828508174909),
new google.maps.LatLng(45.482615296381,9.2038542810421),
new google.maps.LatLng(45.445320738824,9.1555122367407),
new google.maps.LatLng(45.454770238052,9.2259774215776),
new google.maps.LatLng(45.428395471124,9.1522770977499),
new google.maps.LatLng(45.479548736232,9.2104669198347),
new google.maps.LatLng(45.460138194146,9.1906854983623),
new google.maps.LatLng(45.438274463092,9.1798210090832),
new google.maps.LatLng(45.507394801084,9.1254968180551),
new google.maps.LatLng(45.497041055802,9.1063102258693),
new google.maps.LatLng(45.428084278387,9.1525876099536),
new google.maps.LatLng(45.450179710637,9.2164146805132),
new google.maps.LatLng(45.457546162079,9.2062148392008),
new google.maps.LatLng(45.477089862542,9.1685121006692),
new google.maps.LatLng(45.46432520852,9.1805403019656),
new google.maps.LatLng(45.450205473096,9.2162728731521),
new google.maps.LatLng(45.486383493671,9.1601176094402),
new google.maps.LatLng(45.464111973362,9.129648883464),
new google.maps.LatLng(45.460030398038,9.1568395781713),
new google.maps.LatLng(45.45471994274,9.2269996278862),
new google.maps.LatLng(45.508800249178,9.1330320907246),
new google.maps.LatLng(45.502597617023,9.2102634346546),
new google.maps.LatLng(45.491553053439,9.1467088703219),
new google.maps.LatLng(45.47643003457,9.1816503743219),
new google.maps.LatLng(45.457602217812,9.2081471254118),
new google.maps.LatLng(45.528109851553,9.210004486176),
new google.maps.LatLng(45.484911764081,9.1396720300888),
new google.maps.LatLng(45.445362660146,9.1778565110594),
new google.maps.LatLng(45.460445252923,9.1785446587045),
new google.maps.LatLng(45.48591347414,9.1678337074811),
new google.maps.LatLng(45.480476288627,9.2299681451303),
new google.maps.LatLng(45.489413219588,9.2021213168659),
new google.maps.LatLng(45.45419557378,9.1909826103744),
new google.maps.LatLng(45.507484780105,9.1255857892786),
new google.maps.LatLng(45.467731178421,9.1548905074319),
new google.maps.LatLng(45.498351761051,9.1996946211701),
new google.maps.LatLng(45.480765801408,9.1826064888101),
new google.maps.LatLng(45.466385469874,9.2073967482766),
new google.maps.LatLng(45.474692714042,9.2215038648965),
new google.maps.LatLng(45.455165925249,9.1601206145568),
new google.maps.LatLng(45.48513363364,9.1864745941832),
new google.maps.LatLng(45.485559648524,9.1183744961194),
new google.maps.LatLng(45.453167502042,9.1752523149984),
new google.maps.LatLng(45.50247866658,9.1946450723238),
new google.maps.LatLng(45.483502624132,9.2374496285648),
new google.maps.LatLng(45.482745069739,9.2041214427492),
new google.maps.LatLng(45.507364462711,9.1255397962849),
new google.maps.LatLng(45.484538718944,9.208754927841),
new google.maps.LatLng(45.49785883008,9.1067013916072),
new google.maps.LatLng(45.50219989979,9.1946525051785),
new google.maps.LatLng(45.485069509789,9.18653417116),
new google.maps.LatLng(45.488040436692,9.2191394431471),
new google.maps.LatLng(45.455247712198,9.1598939154178),
new google.maps.LatLng(45.502232901447,9.1969196730343),
new google.maps.LatLng(45.472392110487,9.1726812943235),
new google.maps.LatLng(45.507344716829,9.1255662018377),
new google.maps.LatLng(45.459677462992,9.1571197338585),
new google.maps.LatLng(45.480143468569,9.1748754131048),
new google.maps.LatLng(45.455413173359,9.1600531214054),
new google.maps.LatLng(45.471269723852,9.1541527293846)  ] }); 

twistMap.init();


  // ------------------------------------
  // Scroll
  // ------------------------------------
  function scrollToAnchor(id){
      var dest = $(id);
      $('html,body').animate({
      scrollTop: dest.offset().top
    },'slow', function () {
          window.location.hash = id;
      });
  }

  $(".js-scroll").click(function(e) {
      e.preventDefault();
      var id = $(this).attr('href');
      scrollToAnchor(id);
  });

})();
