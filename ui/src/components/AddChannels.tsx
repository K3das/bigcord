import {APIClient, Guild} from "../utils/api";
import {createResource, createSignal, Show} from "solid-js";
import GuildChannels from "./GuildChannels";
import Modal from "./Modal";

const AddChannels = (props: { api: APIClient }) => {
    const [showPopup, setShowPopup] = createSignal<boolean>(false);
    
    const [guildIdInput, setGuildIdInput] = createSignal<string | null>();

    const [guildId, setGuildId] = createSignal<string | null>();
    const [guild, setGuild] = createResource<Guild, string>(guildId, (id) => props.api.getDiscordGuild(id))

    const [selectedChannels, setSelectedChannels] = createSignal<Set<string>>(new Set());

    
    const close = () => {
        setGuildId()
        setGuild.mutate()
        setShowPopup(false)
    }
    
    const done = () => {
        if (selectedChannels().size === 0) {
            return close()
        }
        
        props.api.addChannels(selectedChannels()).then(console.log).catch(console.log)

        close()
    }
    
    return (
        <div class={"my-3"}>
            <button type="button" class={"focus:outline-none text-white bg-green-700 hover:bg-green-800 focus:ring-4 focus:ring-green-300 font-medium text-sm px-5 py-2.5 mr-2 mb-2 dark:bg-green-600 dark:hover:bg-green-700 dark:focus:ring-green-800"} onClick={() => setShowPopup(true)}>Add channels</button>
            
            <Show when={showPopup()}>
                <Modal title={"Add Channels"} close={close} done={done}>
                    <div class="relative p-6 flex-auto w-96">
                        <form onSubmit={(e) => {
                            e.preventDefault()
                            setGuildId(guildIdInput())
                        }}>
                            <div class="relative">
                                <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                                    <svg class="w-4 h-4 text-gray-500 dark:text-gray-400" aria-hidden="true"
                                         xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20">
                                        <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                              d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"/>
                                    </svg>
                                </div>
                                <input type="search" id="default-search"
                                       class={"block w-full p-4 pl-10 text-sm text-gray-900 border border-gray-300 bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" + (guild.error ? " border-red-500" : "")}
                                       placeholder="Search guild by ID" onInput={(e) => setGuildIdInput(e.currentTarget.value)}
                                       required></input>
                                <button type="submit"
                                        class="text-white absolute right-2.5 bottom-2.5 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium text-sm p-2 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
                                        disabled={guild.loading && !guild.error}>{guild.loading && !guild.error ? "⏳" : "🔍"}</button>
                            </div>
                        </form>

                        <Show when={!guild.error && guild()}>
                            <div>
                                <h3 class={"font-bold my-4 text-2xl"}>{guild().name}</h3>
                                <GuildChannels channels={guild().channels} onUpdate={setSelectedChannels}></GuildChannels>
                            </div>
                        </Show>
                    </div>
                </Modal>
            </Show>
        </div>

    )
}

export default AddChannels