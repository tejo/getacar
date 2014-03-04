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


getACar = {
  start : ->
    overLayManager = new OverlayManager();
};

# Dom Ready
$ ->
  getACar.start()
  
  
  
  