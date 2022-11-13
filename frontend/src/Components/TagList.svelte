<script>
    import {createEventDispatcher} from "svelte";

    export let allTags = {}, totalSnippets = 0;
    let selectedTag = '';

    const dispatch = createEventDispatcher();
    function selectTag(name) {
        selectedTag = name;
        dispatch('message', { tag: name });
    }
</script>

{#if totalSnippets}
    <section>
        <div on:click={() => selectTag('')} class:selected={selectedTag === ''}>
            <span>All notes</span>
            <span>{totalSnippets}</span>
        </div>
        {#if allTags}
            {#each Object.entries(allTags) as [tag, count]}
                <div on:click={() => selectTag(tag)} class:selected={selectedTag === tag}>
                    <span>{tag}</span>
                    <span>{count}</span>
                </div>
            {/each}
        {/if}
    </section>
{/if}

<style>
    section {
        flex: 1 1 auto;
        min-width: 120px;
        padding: 1rem;
        border-right: 1px solid rgba(72, 99, 90, 30%);
        overflow-y: auto;
    }
    div {
        display: flex;
        margin-bottom: .75rem;
        cursor: pointer;
        border-bottom: 1px solid transparent;
        padding-bottom: 0.25rem;
    }
    div:last-of-type {
        margin-bottom: 0;
    }
    div:hover {
        color: #D6D0A5;
        border-bottom: 1px dashed rgba(214, 208, 165, 50%);
    }
    span {
        flex: 1;
    }
    div > :first-child {
        display: inline-flex;
        justify-content: flex-start;
    }
    div > :last-child {
        display: inline-flex;
        justify-content: flex-end;
    }
    .selected {
        color: #96A672 !important;
        border-bottom: 1px solid #96A672 !important;
    }
</style>
