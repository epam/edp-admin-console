$(function () {

    !function () {
        $('input.promote-checkbox').prop('checked', false);

        $.each($('.cd-pipeline-applications input'), function () {
            let $appBodyEl = $('.applications .card-body'),
                $inputAppEl = $appBodyEl.find('input[value="' + $(this).data('app-name') + '"]');

            if ($inputAppEl.length !== 0) {
                $inputAppEl.prop('checked', true);

                let $selectDockerStreamEl = $appBodyEl.find('select[name="' + $(this).data('app-name') + '"]');

                if ($selectDockerStreamEl) {
                    $selectDockerStreamEl.val($(this).data('docker-stream-name'));
                    $selectDockerStreamEl.prop('disabled', false);
                }
            }

        });

        $.each($('.applications-to-promote input'), function () {
            let appToPromote = $(this).data('app-name');
            $.each($('.applications .card-body').find('input.promote-checkbox'), function () {
                if ($(this).attr('id') === appToPromote + '-promote') {
                    $(this).prop('checked', true);
                    $(this).prop('disabled', false);
                }
            });
        });

        $.each($('#collapseTwo .card-body .row'), function () {
            if ($(this).find('.app-block input').is(':checked')) {
                $(this).find('.promote-block input').prop('disabled', false);
            }
        });
    }();

    $('.application-checkbox :checkbox').change(function () {
        let $selectEl = $('.select-' + $(this).attr('id')),
            $checkboxEl = $('.checkbox-' + $(this).attr('id'));
        if ($(this).is(':checked')) {
            $selectEl.prop('disabled', false);
            $checkboxEl.prop('disabled', false);
            $('.app-checkbox-error').hide();
        } else {
            $selectEl.prop('disabled', true);
            $checkboxEl.prop('disabled', true);
        }
    });

    $('.update-cd-pipeline').click(function (e) {
        e.preventDefault();

        if (!$('.app-checkbox').is(':checked')) {
            $('.app-checkbox-error').show();
        } else {
            $('.app-checkbox-error').hide();
            $('#updateCDCR').submit();
        }

    });

});