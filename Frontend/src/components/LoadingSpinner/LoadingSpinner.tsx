import "@/components/LoadingSpinner/LoadingSpinner.scss"

export function LoadingSpinner() {
    return (
        <div className="lds-ring">
            <div></div>
            <div></div>
            <div></div>
            <div></div>
        </div>
    );
}