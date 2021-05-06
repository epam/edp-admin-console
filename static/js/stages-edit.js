function editStages() {
    $('#edit-stages').attr('disabled', true);
    $('#save-stages').attr('disabled', false);
    $('.trigger-type-edit').removeClass('hide-element');
    $('.trigger-type-view').addClass('hide-element');
}

function saveStages() {
    $('#edit-stages').attr('disabled', false);
    $('#save-stages').attr('disabled', true);
    $('.trigger-type-edit').addClass('hide-element');
    $('.trigger-type-view').removeClass('hide-element');
    $('#editCDStages').submit();
    showNotification(true, null, 'The pipeline stages has been edited successfully.', 'Editing has been failed.');
}