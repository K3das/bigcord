import {JSX, Show} from "solid-js";
import Title from "./Title";
import {Portal} from "solid-js/web";

const Modal = (props: { children: JSX.Element, title: string, close: () => void, done?: () => void }) => {
    return (
        <Portal>
            <div
                class="justify-center items-center flex overflow-x-hidden overflow-y-auto fixed inset-0 z-50 outline-none focus:outline-none"
            >
                <div class="relative w-auto my-6 mx-auto max-w-3xl max-h-screen">
                    <div
                        class="border-0 shadow-lg relative flex flex-col w-full bg-white outline-none focus:outline-none">
                        <div
                            class="flex items-start justify-between p-5 border-b border-solid border-slate-200 rounded-t">
                            <Title>{props.title}</Title>
                        </div>

                        {props.children}

                        <div
                            class="flex items-center justify-between p-6 border-t border-solid border-slate-200 rounded-b">
                            <button
                                class="text-red-500 background-transparent font-bold uppercase px-6 py-2 text-sm outline-none focus:outline-none mr-1 mb-1 ease-linear transition-all duration-150"
                                type="button"
                                onClick={props.close}
                            >
                                Close
                            </button>
                            <Show when={props.done}>
                                <button
                                    class="text-green-500 background-transparent font-bold uppercase px-6 py-2 text-sm outline-none focus:outline-none mr-1 mb-1 ease-linear transition-all duration-150"
                                    type="button"
                                    onClick={props.done}
                                >
                                    Done
                                </button>
                            </Show>
                        </div>

                    </div>
                </div>
            </div>
            <div class="opacity-25 fixed inset-0 z-40 bg-black"></div>
        </Portal>
    )
}

export default Modal