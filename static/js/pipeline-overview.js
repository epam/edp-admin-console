$(function () {
    $('.tooltip-icon').tooltip();

    !function () {
        $.each($('.applications-to-promote input'), function () {
            let appToPromote = $(this).data('app-name');
            $.each($('.applications-info .edp-table tbody tr'), function () {
                let $promoteEl = $(this).find('.promoteCDPipeline');
                if ($(this).find('.codebaseName').text().trim() === appToPromote) {
                    $promoteEl.find('.promote-checkbox-overview').addClass('promoted').show();
                } else {
                    $promoteEl.find('.promote-checkbox-overview').addClass('cancel').show();
                }
            });

        });
    }();

});
