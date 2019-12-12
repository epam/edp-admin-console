$(function () {
    $('.tooltip-icon').tooltip();

    !function () {
        $.each($('.applications-to-promote input'), function () {
            let appToPromote = $(this).data('app-name');
            $.each($('.applications-info .edp-table tbody tr'), function () {
                let $promoteEl = $(this).find('.promoteCDPipeline');
                if ($(this).find('.codebaseName').text().trim() === appToPromote) {
                    $promoteEl.find('.promote-checkbox-overview').removeClass('cancel').addClass('promoted');
                }
            });

        });
    }();

    !function () {
        $.each($('.platform-link a'), function () {
            if (!$(this).attr('href')) {
                $(this).addClass('hover-popup')
                    .attr('disabled', true)
                    .css('color', '#aaa');
            }
        })
    }();

    let moveLeft = 20,
        moveDown = 10,
        $link = $('a.hover-popup');

    $link.hover(function () {
        $('div#kubernetes-component').show();
    }, function () {
        $('div#kubernetes-component').hide();
    });

    $link.mousemove(function (e) {
        $("div#kubernetes-component")
            .css('top', e.pageY + moveDown).css('left', e.pageX + moveLeft);
    });

    $('.platform-link a.edp-link').click(function (e) {
        if ($(this).attr('disabled')) {
            e.preventDefault();
        }
    });

});
