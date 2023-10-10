import {createResource, For, Match, onCleanup, onMount, Switch} from "solid-js";
import {APIClient, ChannelType, GetChannelsResponse, JobState, State} from "../utils/api";

const ChannelStates = (props: { api: APIClient }) => {
    const [channels, setChannels] = createResource<GetChannelsResponse>(() => props.api.getChannels())

    onMount(() => {
        let interval = setInterval(() => {
            if (channels.error == undefined && !channels()) {
                return
            }
            setChannels.refetch()
        }, 1000)
        
        onCleanup(() => clearInterval(interval))
    })

    return (
        <div class={"overflow-auto my-3"}>
            <Switch fallback={<pre>Error :(</pre>}>
                <Match when={channels.error == undefined && channels()}>
                    <Table data={channels()} />
                </Match>
                <Match when={channels.loading}>
                    <pre>Loading...</pre>
                </Match>
                <Match when={channels.error}>
                        Error while fetching channels:
                        <pre>{channels.error.toString()}</pre>
                </Match>
            </Switch>
        </div>
    );
};

const Table = (props: { data: GetChannelsResponse }) => {
    return (
        <table class={"min-w-full divide-y divide-gray-200"}>
            <thead>
            <tr>
                <th class={"border p-2"}>ID</th>
                <th class={"border p-2"}>Guild ID</th>
                <th class={"border p-2"}>Type</th>
                <th class={"border p-2"}>Name</th>
                <th class={"border p-2"}>Offset</th>
                <th class={"border p-2"}>State</th>
                <th class={"border p-2"}>Job State</th>
            </tr>
            </thead>
            <tbody>
                <For each={props.data.channels}>{(channel) => (
                    <tr>
                        <td class={"border p-2"}>{channel.id}</td>
                        <td class={"border p-2"}>{channel.guild_id}</td>
                        <td class={"border p-2"}>{ChannelType[channel.type]}</td>
                        <td class={"border p-2"}>{channel.name}</td>
                        <td class={"border p-2"}>{channel.offset}</td>
                        <td class={"border p-2"}>{State[channel.state]}</td>
                        <td class={"border p-2"}>{JobState[channel.job_state]}</td>
                    </tr>
                )}</For>
            </tbody>
        </table>
    );
};

export default ChannelStates;
