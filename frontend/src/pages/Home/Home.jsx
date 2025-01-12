import Header from "../../components/Header/Header";
import CustomGrid from "../../components/CustomGrid/CustomGrid";
import styles from "./Home.module.css";
const Home = () => {
  return (
    <>
      <div class={styles.container}>
        <Header />
        <CustomGrid />
      </div>
    </>
  );
};

export default Home;
