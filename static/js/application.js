$(function () {

    $(document).ready(function () {
        let $successPopupEl = $('#successPopup');
        let displayPopup = $successPopupEl.data("display");
        if (displayPopup) {
            $successPopupEl.modal('show');
        }
    });

});