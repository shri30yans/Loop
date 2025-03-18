const backgroundImages = [
    "https://www.liquidplanner.com/wp-content/uploads/2019/04/HiRes-17.jpg",
    "https://images.unsplash.com/photo-1522071820081-009f0129c71c",
    "https://images.unsplash.com/photo-1507238691740-187a5b1d37b8",
    "https://images.unsplash.com/photo-1531498860502-7c67cf02f657",
    "https://images.unsplash.com/photo-1556761175-4b46a572b786",
    "https://boingboing.net/wp-content/uploads/2013/10/FY4TBHSHMMFBB4V.LARGE_2.jpg",
    "https://content.instructables.com/ORIG/FOU/QRZL/KBKSQ17D/FOUQRZLKBKSQ17D.jpg"
    ,"https://cdn.instructables.com/ORIG/FT4/67OF/I6W2S2GW/FT467OFI6W2S2GW.jpg?frame=1",
    "https://content.instructables.com/ORIG/FBJ/BJ4I/K14TZFFK/FBJBJ4IK14TZFFK.jpg?auto=webp&frame=1&width=2100"
  ];
  
  export const getRandomBackground = () => {
      const randomIndex = Math.floor(Math.random() * backgroundImages.length);
      return backgroundImages[randomIndex];
    };