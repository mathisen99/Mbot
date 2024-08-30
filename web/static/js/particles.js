particlesJS("particles-js", {
  "particles": {
    "number": {
      "value": 70, // Slightly increase the number of particles
      "density": {
        "enable": true,
        "value_area": 800
      }
    },
    "color": {
      "value": "#7e9aa8" // Subtle blue-gray color to match the accent
    },
    "shape": {
      "type": "circle", // Circles are clean and minimalistic
      "stroke": {
        "width": 0,
        "color": "#000000"
      },
      "polygon": {
        "nb_sides": 6 // Hexagons for a slightly more complex shape
      }
    },
    "opacity": {
      "value": 0.3, // Lower opacity for a more subtle appearance
      "random": true, // Make opacity random to add variation
      "anim": {
        "enable": true,
        "speed": 1,
        "opacity_min": 0.1,
        "sync": false
      }
    },
    "size": {
      "value": 4, // Slightly larger particles
      "random": true, // Random sizes for variation
      "anim": {
        "enable": true,
        "speed": 20, // Slow down the size animation
        "size_min": 0.1,
        "sync": false
      }
    },
    "line_linked": {
      "enable": true,
      "distance": 150,
      "color": "#a3b4c6", // Use a lighter blue-gray for connecting lines
      "opacity": 0.2, // Lower opacity for less intense connections
      "width": 1
    },
    "move": {
      "enable": true,
      "speed": 2, // Slower movement for a more relaxed feel
      "direction": "none",
      "random": true, // Randomize movement to add dynamism
      "straight": false,
      "out_mode": "out",
      "bounce": false,
      "attract": {
        "enable": false,
        "rotateX": 600,
        "rotateY": 1200
      }
    }
  },
  "interactivity": {
    "detect_on": "canvas",
    "events": {
      "onhover": {
        "enable": true,
        "mode": "grab" // Use "grab" mode for a more interactive feel
      },
      "onclick": {
        "enable": true,
        "mode": "push"
      },
      "resize": true
    },
    "modes": {
      "grab": {
        "distance": 300, // Adjust grab distance to be more responsive
        "line_linked": {
          "opacity": 0.4 // Higher opacity on hover for better visibility
        }
      },
      "bubble": {
        "distance": 300,
        "size": 10, // Smaller bubble size for subtle effect
        "duration": 1,
        "opacity": 0.5, // Lower opacity for a more subdued bubble effect
        "speed": 3
      },
      "repulse": {
        "distance": 150, // Reduce repulse distance for subtle interaction
        "duration": 0.4
      },
      "push": {
        "particles_nb": 4
      },
      "remove": {
        "particles_nb": 2
      }
    }
  },
  "retina_detect": true
});
