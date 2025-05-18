import './listen';
import Glide from '@glidejs/glide';
import '@glidejs/glide/dist/css/glide.core.min.css';
import '@glidejs/glide/dist/css/glide.theme.min.css';

new Glide('.glide', {
    type: 'slider',
    focusAt: 'center',
}).mount();