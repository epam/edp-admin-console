$(function () {

    const maxJiraFields = 3;
    const defaultJiraIssueField = 'fixVersions';
    const defaultJiraPattern = 'EDP_VERSION-EDP_COMPONENT';

    $($('section select.jiraIssueFields')[0]).selectize({
        create: true,
        sortField: {
            field: 'text',
            direction: 'asc'
        },
        dropdownParent: 'body'
    });

    !function preloadJiraMetadataFields() {
        $('select.jiraServer option').each(function (i, v) {
            let jiraServer = $(v).val();
            let fields = {
                "components": "Component/s",
                "labels": "Labels",
                "fixVersions": "Fix Version/s"
            };
            _setData(jiraServer, fields);
        });
    }();

    function _setData(jiraServer, fields) {
        let $advSettingsBlock = $('div.advanced-settings-block'),
            confData = $advSettingsBlock.attr('data-config');
        if (!!confData) {
            let tempConfData = $.parseJSON(confData);
            tempConfData[jiraServer] = fields;
            $advSettingsBlock.attr('data-config', JSON.stringify(tempConfData))
            return
        }
        $advSettingsBlock.attr('data-config', JSON.stringify({[jiraServer]: fields}))
    }

    $('#jiraServerToggle').change(function () {
        if ($(this).is(':checked')) {
            let jiraServer = $('select.jiraServer option:selected').val(),
                fields = $.parseJSON($('div.advanced-settings-block').attr('data-config'))[jiraServer];
            $.each(fields, function (i, item) {
                $('select.jiraIssueFields')[0].selectize.addOption({value: i, text: item});
            });
            _setDefaultValues();
            return
        }
    });

    function _setDefaultValues() {
        $('section select.jiraIssueFields')[0].selectize.setValue(defaultJiraIssueField);
        $('section div.jira-issue-metadata-row input.jiraPattern').val(defaultJiraPattern);
    }

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
        let jiraServer = $('select.jiraServer option:selected').val(),
            fields = $.parseJSON($('div.advanced-settings-block').attr('data-config'))[jiraServer];
        $.each(fields, function (i, item) {
            if ($.inArray(i, selectedValues) !== -1) {
                return;
            }
            $newSelect[$newSelect.length - 1].selectize.addOption({value: i, text: item})
        });

        tryToDisableButton();
    });

    $('.delete.remove-jira-issue-metadata-row').click(function () {
        removeJiraIssueRow.bind(this)();
    });

    $('.jiraServer').change(function () {
        $('section div.jira-issue-metadata-row').remove();
        enableButton();
    });

    $('.jiraIssueFields.jiraFieldName').change(function () {
        toggleSelectOptions.bind(this)();
    });

});