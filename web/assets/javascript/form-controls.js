const FormControls = (element) => {

    //event delegation bound to all data-control="select-submit"
    document.addEventListener("change", function (e) {
        if (e.target?.dataset?.control == 'select-submit') {
            e.target.closest('form')?.submit();
        }
    })

}

export default FormControls;