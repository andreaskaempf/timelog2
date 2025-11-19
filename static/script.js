// Handler to confirm deletion of project
function confirmProjectDeletion(id) {
    if ( confirm('Are you sure you want to delete this project?') ) {
        window.location.href = '/delete_project/' + id;
    }
}

// Handler to confirm deletion of contact
function confirmContactDeletion(id) {
    if ( confirm('Are you sure you want to delete this contact?') ) {
        window.location.href = '/delete_contact/' + id;
    }
}


// Handler to confirm deletion of work entry
function confirmWorkDeletion(id) {
    if ( confirm('Are you sure you want to delete this entry?') ) {
        window.location.href = '/delete_work/' + id;
    }
}

// Confirm deletion of project/contact link
function confirmProjectLinkDeletion(contactId, projectId, projectName) {
    if ( confirm('Remove link to project "' + projectName + '"?') ) {
        window.location.href = '/del_contact_project?cid=' + contactId + '&pid=' + projectId;
    }
}
