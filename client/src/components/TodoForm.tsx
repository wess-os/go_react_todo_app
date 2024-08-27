import { Button, Flex, Input, Spinner, Text } from "@chakra-ui/react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { BASE_URL } from "../App";

const TodoForm = () => {
    const [newTodo, setNewTodo] = useState("");

    const queryClient = useQueryClient();

    const {mutate: createTodo, isPending: isCreating} = useMutation({
        mutationKey:["createTodo"],
        mutationFn: async (e:React.FormEvent) => {
            e.preventDefault();
            try {
                const res = await fetch(BASE_URL + "todos", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        body: newTodo,
                    }),
                });

                const data = await res.json();

                if (!res.ok) {
                    throw new Error(data.error || "Something went wrong");
                }

                setNewTodo("");

                return data;
            } catch (error:any) {
                console.log(error);
            }
        },
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey: ["todos"]});
        },
        onError: (error:any) => {
            alert(error.message);
        }
    });

    return (
        <form onSubmit={createTodo}>
            <Text fontSize={"4xl"} textTransform={"uppercase"} fontWeight={"bold"} textAlign={"center"}>
                Tarefas do dia
            </Text>
            <Flex gap={2}>
                <Input
                    type="text"
                    value={newTodo}
                    onChange={(e) => setNewTodo(e.target.value)}
                    ref={(input) => input && input.focus()}
                    borderColor={"gray.400"}
                />
            </Flex>
            <Flex gap={2} alignItems={"center"} justifyContent={"center"} my={2}>
                <Button
                        mx={2}
                        type="submit"
                        _active={{
                            transform: "scale(.97)",
                        }}
                        backgroundColor={"blue.400"}
                    >
                    {isCreating ? <Spinner size="xs" /> : 'Adicionar'}
                </Button>
            </Flex>
        </form>
    );
}

export default TodoForm