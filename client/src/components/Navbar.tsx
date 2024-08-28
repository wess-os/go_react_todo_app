import { Box, Flex, Button, useColorModeValue, useColorMode, Container } from "@chakra-ui/react";
import { IoMoon } from "react-icons/io5";
import { LuSun } from "react-icons/lu";
import { IoMdLogOut } from "react-icons/io";
import { FaUserCircle } from "react-icons/fa";

export default function Navbar() {
    const { colorMode, toggleColorMode } = useColorMode();

    const handleLogout = () => {
        localStorage.removeItem('token');

        window.location.href = '/login';
    };

    return (
        <Container maxW={"900px"}>
            <Box bg={useColorModeValue("gray.200", "gray.700")} px={4} my={4} borderRadius={"5"}>
                <Flex h={16} alignItems={"center"} justifyContent={"space-between"}>
                        <Flex
                            gap={2}
                            textAlign={"center"}
                            alignItems={"center"}
                            justifyContent={"center"}
                        >
                            <FaUserCircle size={30} />
                        </Flex>
                    <Flex
                        justifyContent={"center"}
                        alignItems={"center"}
                        gap={3}
                        display={{ base: "none", sm: "flex" }}
                    >
                        <h1>Controle de Tarefas</h1>
                    </Flex>

                    <Flex alignItems={"center"} gap={3}>
                        <Button onClick={toggleColorMode}>
                            {colorMode === "light" ? <IoMoon /> : <LuSun size={20} />}
                        </Button>
                        <Button onClick={handleLogout}>
                            <IoMdLogOut size={20} />
                        </Button>
                    </Flex>
                </Flex>
            </Box>
        </Container>
    );
}