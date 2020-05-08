$(function () {

    $('.update-codebase').click(function () {
        if (_arePatternsValid()) {
            $('#updateCodebase').submit();
        }
    });

    function _arePatternsValid() {
        let $commitMsgEl = $('#commitMessagePattern'),
            $ticketNumberEl = $('#ticketNamePattern'),
            isCommitMessageRegexValid = $commitMsgEl.val().length !== 0,
            isTicketNameRegexValid = $ticketNumberEl.val().length !== 0;

        $('.invalid-feedback.commitMessagePattern').hide();
        $commitMsgEl.removeClass('is-invalid');
        $('.invalid-feedback.ticketNamePattern').hide();
        $ticketNumberEl.removeClass('is-invalid');

        if (!isCommitMessageRegexValid) {
            $('.invalid-feedback.commitMessagePattern').show();
            $commitMsgEl.addClass('is-invalid');
        }

        if (!isTicketNameRegexValid) {
            $('.invalid-feedback.ticketNamePattern').show();
            $ticketNumberEl.addClass('is-invalid');
        }

        return isCommitMessageRegexValid && isTicketNameRegexValid;
    }

});