class ButtonsManager
  constructor: () ->
    #$('form[name=geoCodeForm]').on('click','.bt-search', (e) - )

class OverlayManager
  constructor: () ->
    @container = $("div.container")
    @triggerBttn = $("#trigger-overlay")
    @overlay = $("div.overlay")
    @closeBttn = @overlay.find("button.overlay-close")
    @transEndEventNames =
      WebkitTransition: "webkitTransitionEnd"
      MozTransition: "transitionend"
      OTransition: "oTransitionEnd"
      msTransition: "MSTransitionEnd"
      transition: "transitionend"

    @transEndEventName = @transEndEventNames[Modernizr.prefixed("transition")]
    @support = transitions: Modernizr.csstransitions

    @triggerBttn.on "click", @toggleOverlay
    @closeBttn.on "click", @toggleOverlay

  toggleOverlay : (e) =>
    e.preventDefault();
    _clicked  = $(this)
    if @overlay.hasClass("open")
      @overlay.removeClass "open"
      @container.removeClass "overlay-open"
      @overlay.addClass "close"
      onEndTransitionFn = (ev) =>
        if @support.transitions
          if ev.propertyName isnt "visibility" 
            return
          _clicked.off @transEndEventName, onEndTransitionFn
        @overlay.removeClass "close"
        return
      if @support.transitions
        @overlay.on @transEndEventName, onEndTransitionFn
      else
        onEndTransitionFn()
    else unless @overlay.hasClass("close")
      @overlay.addClass "open"
      @container.addClass "overlay-open"
    return



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


getACar = {
  start : ->
    #btsManager = new ButtonsManager();
    overLayManager = new OverlayManager();
};

# Dom Ready
$ ->
  getACar.start()
  
  
  
  