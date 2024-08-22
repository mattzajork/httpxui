<script>
    import { onMount, onDestroy } from 'svelte';
    import bgimage from '$lib/assets/imgbg.png';

    export let src;
    export let fullView;

    let imgElement;
    let observer;

    function handleIntersect(entries) {
        entries.forEach((entry) => {
            if (entry.isIntersecting) {
                if (entry.target.dataset.src) {
                    entry.target.src = entry.target.dataset.src;
                }
            } else {
                entry.target.src = bgimage;
            }
        });
    }

    onMount(() => {
        let options = {
            root: null,
            rootMargin: "0px",
            threshold: .1 
        };
        observer = new IntersectionObserver(handleIntersect, options);
        observer.observe(imgElement);
    });

    onDestroy(() => {
        if (observer) {
            observer.unobserve(imgElement);
        }
    });
</script>

<img
    bind:this={imgElement}
    class="{!fullView ? 'object-cover object-top aspect-video' : ''} border border-pdgrayoutline" 
    data-src={src}
    src={src}
    on:error={() => imgElement.src = bgimage}
    alt=""
/>
