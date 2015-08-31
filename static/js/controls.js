(function($) {
    $.fn.tableHover = function () {
        var table = $(this);

        var rows = table.find('tr:gt(0)');
        for (var index = 0; index < rows.length; index++) {
            rows.eq(index).mouseover(function () {
                $(this).addClass("rowover");
            }).mouseout(function () {
                $(this).removeClass("rowover");
            });
        }
        
    };
})(jQuery);