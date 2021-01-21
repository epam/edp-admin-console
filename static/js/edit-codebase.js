$(function () {

    $('.update-codebase').click(function () {
        if (_isJiraConfigValid()) {
            $('#updateCodebase').submit();
        }
    });

    function _isJiraConfigValid() {
        let $commitMsgEl = $('#commitMessagePattern'),
            $ticketNumberEl = $('#ticketNamePattern'),
            isCommitMessageRegexValid = $commitMsgEl.val().length !== 0,
            isTicketNameRegexValid = $ticketNumberEl.val().length !== 0,
            $jiraInputs = $('section .jira-issue-metadata-row input.jiraPattern'),
            $jiraSelects = $('section .jira-issue-metadata-row select.jiraIssueFields'),
            areJiraInputFieldsValid = true,
            areJiraSelectsValid = true;

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

        if ($('div.jiraIssueMetadata div.jira-issue-metadata-row').length > 0) {
            $.each($jiraInputs, function () {
                if ($(this).val() === "") {
                    areJiraInputFieldsValid = false;
                    $(this)
                        .parents('.jira-issue-metadata-row')
                        .parents('.jiraIssueMetadata')
                        .find('.invalid-feedback.jira-row-invalid-msg')
                        .show();

                    $(this).addClass('is-invalid');
                } else {
                    areJiraInputFieldsValid = true;
                    $(this).removeClass('is-invalid');
                }
            });

            $.each($jiraSelects, function () {
                if ($(this).find('option:selected').val() === "") {
                    areJiraSelectsValid = false;
                    $(this)
                        .parents('.jira-issue-metadata-row')
                        .parents('.jiraIssueMetadata')
                        .find('.invalid-feedback.jira-row-invalid-msg')
                        .show();

                    $(this).addClass('is-invalid');
                } else {
                    areJiraSelectsValid = true;
                    $(this).removeClass('is-invalid');
                }
            });

            if (areJiraInputFieldsValid && areJiraSelectsValid) {
                $('div.jira-row-invalid-msg').hide();
            }
        }

        return isCommitMessageRegexValid && isTicketNameRegexValid && areJiraInputFieldsValid && areJiraSelectsValid;
    }

});