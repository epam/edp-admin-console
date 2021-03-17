$(function () {

    const maxJiraFields = 3;
    let fields = {
        "components": "Component/s",
        "labels": "Labels",
        "fixVersions": "Fix Version/s"
    };

    !function toggleJiraViewOnInit() {
        toggleJiraView.bind($('#jiraServerToggle'))();
    }();

    !function preloadJiraMetadataFields() {
        let $jiraMetadataBlock = $('div .jiraIssueMetadata'),
            jiraFields = $jiraMetadataBlock.attr('data-conf');
        if (!jiraFields) {
            $jiraMetadataBlock.attr('data-conf', JSON.stringify(fields))
            return
        }

        $jiraMetadataBlock.removeClass('hide-element');

        let jiraIssueMetadataFields = JSON.parse(jiraFields);
        $.each(jiraIssueMetadataFields, function (k, v) {
            let $row = getTemplate();
            $row.insertBefore($('.add-jira-field'));

            let removeButtonId = generateId();
            $row.find('.remove-jira-issue-metadata-row').addClass(removeButtonId)

            let selectId = generateId();
            $row.find('select').addClass(selectId)

            let $newSelect = $('section select.jiraIssueFields').last()
            $($newSelect).selectize({
                create: true,
                sortField: {
                    field: 'text',
                    direction: 'asc'
                },
                dropdownParent: 'body'
            });
            $row
                .removeClass('full-template')
                .removeClass('partial-template')
                .removeClass('hide-element');

            let $deleteBtn = $newSelect.parents('.jira-issue-metadata-row').find(`button.remove-jira-issue-metadata-row.${removeButtonId}`);
            $('.jiraIssueMetadata').on('click', `.remove-jira-issue-metadata-row.${removeButtonId}`, removeJiraIssueRow.bind($deleteBtn));

            let $select = $newSelect.parents('.jira-issue-metadata-row').find(`select.${selectId}`);
            $('.jiraIssueMetadata').on('change', `.${selectId}`, toggleSelectOption.bind($select));

            let selectedValues = $('section select.jiraIssueFields').map(function () {
                return this.value
            }).get();
            $.each(fields, function (i, v) {
                if ($.inArray(i, selectedValues) !== -1) {
                    return;
                }
                $newSelect[$newSelect.length - 1].selectize.addOption({value: i, text: v})
            });

            $newSelect[0].selectize.setValue(k);
            $newSelect.parents('div .jira-issue-metadata-row').find('input.jiraPattern').val(v);

            tryToDisableButton();
        });
    }();

    $('.add-jira-field').click(function () {
        if ($('.jiraIssueMetadata .jira-issue-metadata-row').length === maxJiraFields) {
            return
        }

        let $row = getTemplate();
        $row.insertBefore($('.add-jira-field'));

        let removeButtonId = generateId();
        $row.find('.remove-jira-issue-metadata-row').addClass(removeButtonId)

        let selectId = generateId();
        $row.find('select').addClass(selectId)

        let $newSelect = $('section select.jiraIssueFields').last()
        $($newSelect).selectize({
            create: true,
            sortField: {
                field: 'text',
                direction: 'asc'
            },
            dropdownParent: 'body'
        });
        $row
            .removeClass('full-template')
            .removeClass('partial-template')
            .removeClass('hide-element');

        let $deleteBtn = $newSelect.parents('.jira-issue-metadata-row').find(`button.remove-jira-issue-metadata-row.${removeButtonId}`);
        $('.jiraIssueMetadata').on('click', `.remove-jira-issue-metadata-row.${removeButtonId}`, removeJiraIssueRow.bind($deleteBtn));

        let $select = $newSelect.parents('.jira-issue-metadata-row').find(`select.${selectId}`);
        $('.jiraIssueMetadata').on('change', `.${selectId}`, toggleSelectOption.bind($select));

        let selectedValues = $('section select.jiraIssueFields').map(function () {
            return this.value
        }).get();
        $.each(fields, function (i, v) {
            if ($.inArray(i, selectedValues) !== -1) {
                return;
            }
            $newSelect[$newSelect.length - 1].selectize.addOption({value: i, text: v})
        });

        tryToDisableButton();
    });


    $('#jiraServerToggle').change(function () {
        toggleJiraView.bind($(this))();
    });

});