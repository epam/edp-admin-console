function deleteConfirmation() {
    let $modal = $("#delete-confirmation"),
        targetName = $modal.data('name'),
        confirmationName = $modal.find('#entity-name').val(),
        $errName = $modal.find('.invalid-feedback.different-name'),
        $errUsed = $modal.find('.invalid-feedback.server-error');
    if (targetName !== confirmationName) {
        $errName.show();
        return
    }
    $errName.hide();
    $errUsed.hide();
    $errUsed.text('');
    $("#delete-action").submit();
}

function closeConfirmation() {
    let $modal = $("#delete-confirmation");
    $modal.find('.invalid-feedback.different-name').hide();
    $modal.find('.invalid-feedback.server-error').text('').hide();
    $modal.find('#entity-name').val('');
}