class MapManager
  constructor: (@map_canvas) ->
    @map_center = new google.maps.LatLng(-34.397, 150.644)
    @map_zoom = 8
    google.maps.event.addDomListener(window, 'load', @init)

  init: () ->
    console.log 'load'
    console.log $('#map-canvas').length
    mapOptions = {center: @map_center,zoom: @map_zoom}
    @map_handle = new google.maps.Map(@map_canvas, mapOptions)
    
    console.log @map_handle
    
    
  
  
  ###
  loadScript: () ->
    script = document.createElement('script')
    script.type = 'text/javascript'
    script.src = 'https://maps.googleapis.com/maps/api/js?v=3.exp&sensor=false&' +
      'callback=initialize'
    document.body.appendChild(script)
  ###


# Dom Ready
$ ->
  console.log 'ready'
  #mapManager = new MapManager(document.getElementById("map-canvas"))
  