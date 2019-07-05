$(function () {

    !function () {
        $.each($('.cd-pipeline-applications input'), function () {
            let $appBodyEl = $('.applications .card-body'),
                $inputAppEl = $appBodyEl.find('input[value="' + $(this).data('app-name') + '"]');

            if ($inputAppEl.length !== 0) {
                $inputAppEl.prop('checked', true);

                let $selectBranchEl = $appBodyEl.find('select[name="' + $(this).data('app-name') + '"]');

                if ($selectBranchEl) {
                    $selectBranchEl.val($(this).data('branch-name'))
                    $selectBranchEl.prop('disabled', false);
                }
            }

        });
    }();

    $('.application-checkbox :checkbox').change(function () {
        let $selectEl = $('.select-' + $(this).attr('id'));
        if ($(this).is(':checked')) {
            $selectEl.prop('disabled', false);
            $('.app-checkbox-error').hide();
        } else {
            $selectEl.prop('disabled', true);
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