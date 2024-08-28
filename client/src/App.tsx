import { BrowserRouter as Router, Route, Routes, useLocation } from "react-router-dom";
import { Container, Flex, Stack } from "@chakra-ui/react";
import Navbar from "./components/Navbar";
import TodoForm from "./components/TodoForm";
import TodoList from "./components/TodoList";
import Login from "./components/Login";
import Register from "./components/Register";

export const BASE_URL = import.meta.env.MODE === "development" ? "http://localhost:4000/api/" : "/api";

function App() {
  const location = useLocation();

  const showNavbar = location.pathname !== '/login' && location.pathname !== '/register';

  return (
    <Stack h="100vh">
      {showNavbar && <Navbar />}
      <Flex flex="1" justifyContent="center" alignItems="center">
        <Routes>
          <Route path="/" element={
            <Container>
              <TodoForm />
              <TodoList />
            </Container>
          } />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
        </Routes>
      </Flex>
    </Stack>
  );
}

const Main = () => (
  <Router>
    <App />
  </Router>
);

export default Main;
